package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/dto"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var cartTracer = otel.Tracer("bitedash/internal/service/cart")

type CartService struct {
	queries *db.Queries
	db      *sql.DB
}

func NewCartService(q *db.Queries, db *sql.DB) *CartService {
	return &CartService{queries: q, db: db}
}

func (s *CartService) AddItem(ctx context.Context, userID uuid.UUID, req dto.AddCartItemRequest) error {
	ctx, span := cartTracer.Start(ctx, "CartService.AddItem")
	defer span.End()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := s.queries.WithTx(tx)

	product, err := qtx.GetProductByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	if !product.IsAvailable || product.KitchenQuantity < 1 {
		return ErrProductNotAvailable
	}

	cart, err := qtx.GetActiveCartByUser(ctx, userID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		cart, err = qtx.CreateCart(ctx, db.CreateCartParams{
			ID:     uuid.New(),
			UserID: userID,
		})
		if err != nil {
			return err
		}
	}

	var existingQty int32
	currentQty, err := qtx.GetCartItemQuantity(ctx, db.GetCartItemQuantityParams{
		CartID:    cart.ID,
		ProductID: req.ProductID,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
		existingQty = currentQty
	}

	if existingQty+req.Quantity > product.KitchenQuantity {
		return ErrInsufficientStock
	}

	_, err = qtx.UpsertCartItem(ctx, db.UpsertCartItemParams{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CartService) GetCart(ctx context.Context, cartID uuid.UUID) (*dto.CartResponse, error) {
	ctx, span := cartTracer.Start(ctx, "CartService.GetCart")
	defer span.End()

	cart, err := s.queries.GetActiveCartByUser(ctx, cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &dto.CartResponse{
				CartID:      "",
				Items:       []dto.CartItemDto{},
				TotalAmount: 0,
				ItemsCount:  0,
			}, nil
		}
		return nil, err
	}

	rows, err := s.queries.GetCartItemsWithProducts(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	resp := &dto.CartResponse{
		CartID: cart.ID.String(),
		Items:  make([]dto.CartItemDto, 0, len(rows)),
	}

	for _, row := range rows {
		price, err := strconv.ParseFloat(row.Price, 64)
		if err != nil {
			return nil, err
		}

		lineTotal := price * float64(row.Quantity)

		resp.Items = append(resp.Items, dto.CartItemDto{
			ProductID: row.ProductID.String(),
			Name:      row.Name,
			Price:     price,
			Quantity:  row.Quantity,
			LineTotal: lineTotal,
		})
		resp.TotalAmount += lineTotal
		resp.ItemsCount += row.Quantity
	}
	return resp, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	ctx, span := cartTracer.Start(ctx, "CartService.ClearCart")
	defer span.End()

	cart, err := s.queries.GetActiveCartByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCartNotFound
		}
		return err
	}

	return s.queries.ClearCart(ctx, cart.ID)
}

func (s *CartService) RemoveCartItem(ctx context.Context, userID, productID uuid.UUID) error {
	ctx, span := cartTracer.Start(ctx, "CartService.RemoveCartItem")
	defer span.End()

	cart, err := s.queries.GetActiveCartByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCartNotFound
		}
		return err
	}

	rowsAffected, err := s.queries.RemoveCartItem(ctx, db.RemoveCartItemParams{
		CartID:    cart.ID,
		ProductID: productID,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCartItemNotFound
	}

	return nil
}
