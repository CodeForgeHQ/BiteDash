package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/dto"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
)

var orderTracer = otel.Tracer("bitedash/internal/service/order")

type OrderService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewOrderService(dbConn *sql.DB, qb *db.Queries) *OrderService {
	return &OrderService{
		db:      dbConn,
		queries: qb,
	}
}

func (s *OrderService) MakeOrder(
	ctx context.Context,
	userID uuid.UUID,
) (*dto.CheckoutResponse, error) {
	ctx, span := orderTracer.Start(ctx, "OrderService.MakeOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("userID", userID.String()),
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := s.queries.WithTx(tx)

	cart, err := qtx.GetActiveCartByUser(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, "failed to get active cart")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCartNotFound
		}
		return nil, fmt.Errorf("failed to get cart for user: %w", err)
	}

	items, err := qtx.GetCartItemsWithProductStock(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	if len(items) == 0 {
		return nil, ErrCartEmpty
	}

	var totalAmount float64
	var itemsCount int32

	type makeOrderItem struct {
		ProductID   uuid.UUID
		ProductName string
		UnitPrice   float64
		Quantity    int32
		LineTotal   float64
	}

	orderItems := make([]makeOrderItem, 0, len(items))

	for _, item := range items {
		if !item.IsAvailable {
			return nil, fmt.Errorf("%w: %s", ErrProductNotAvailable, item.Name)
		}
		if item.KitchenQuantity < item.Quantity {
			return nil, fmt.Errorf("%w: %s", ErrInsufficientStock, item.Name)
		}

		unitPrice, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price for product %s: %w", item.Name, err)
		}

		lineTotal := unitPrice * float64(item.Quantity)
		totalAmount += lineTotal
		itemsCount += item.Quantity

		orderItems = append(orderItems, makeOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.Name,
			UnitPrice:   unitPrice,
			Quantity:    item.Quantity,
			LineTotal:   lineTotal,
		})
	}

	orderID := uuid.New()

	for _, item := range orderItems {
		rowsAffected, err := qtx.DecreaseProductKitchenQuantity(ctx, db.DecreaseProductKitchenQuantityParams{
			ID:              item.ProductID,
			KitchenQuantity: item.Quantity,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decrease kitchen stock for %s: %w", item.ProductName, err)
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("%w: %s", ErrInsufficientStock, item.ProductName)
		}
	}

	err = qtx.CreateOrder(ctx, db.CreateOrderParams{
		ID:          orderID,
		UserID:      userID,
		Status:      "pending",
		TotalAmount: fmt.Sprintf("%.2f", totalAmount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range orderItems {
		err = qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitPrice:   fmt.Sprintf("%.2f", item.UnitPrice),
			Quantity:    item.Quantity,
			LineTotal:   fmt.Sprintf("%.2f", item.LineTotal),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	err = qtx.DeleteCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &dto.CheckoutResponse{
		OrderID:     orderID.String(),
		Status:      "pending",
		TotalAmount: totalAmount,
		ItemsCount:  itemsCount,
	}, nil
}

func (s *OrderService) ListMyOrders(ctx context.Context, userID uuid.UUID) ([]dto.OrderDTO, error) {
	ctx, span := orderTracer.Start(ctx, "OrderService.ListMyOrders")
	defer span.End()

	orders, err := s.queries.ListOrdersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list orders by user: %w", err)
	}

	result := make([]dto.OrderDTO, 0, len(orders))

	for _, order := range orders {
		orderDTO, err := s.buildOrderDTO(ctx, order)
		if err != nil {
			return nil, err
		}

		result = append(result, *orderDTO)
	}

	return result, nil
}

func (s *OrderService) GetOrderByID(
	ctx context.Context,
	userID uuid.UUID,
	orderID uuid.UUID,
) (*dto.OrderDTO, error) {
	ctx, span := orderTracer.Start(ctx, "OrderService.GetOrderByID")
	defer span.End()

	order, err := s.queries.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf("get order by id: %w", err)
	}

	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	return s.buildOrderDTO(ctx, order)
}

func (s *OrderService) buildOrderDTO(ctx context.Context, order db.Order) (*dto.OrderDTO, error) {
	totalAmount, err := strconv.ParseFloat(order.TotalAmount, 64)
	if err != nil {
		return nil, fmt.Errorf("parse order total amount: %w", err)
	}

	items, err := s.queries.ListOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("list order items: %w", err)
	}

	itemDTOs := make([]dto.OrderDetailsDTO, 0, len(items))

	for _, item := range items {
		unitPrice, err := strconv.ParseFloat(item.UnitPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("parse order item unit price: %w", err)
		}

		lineTotal, err := strconv.ParseFloat(item.LineTotal, 64)
		if err != nil {
			return nil, fmt.Errorf("parse order item line total: %w", err)
		}

		itemDTOs = append(itemDTOs, dto.OrderDetailsDTO{
			ProductID:   item.ProductID.String(),
			ProductName: item.ProductName,
			UnitPrice:   unitPrice,
			Quantity:    item.Quantity,
			LineTotal:   lineTotal,
		})
	}

	createdAt := ""
	if order.CreatedAt.Valid {
		createdAt = order.CreatedAt.Time.Format(time.RFC3339)
	}

	return &dto.OrderDTO{
		OrderID:     order.ID.String(),
		Status:      order.Status,
		TotalAmount: totalAmount,
		Items:       itemDTOs,
		CreatedAt:   createdAt,
	}, nil
}
