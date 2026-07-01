package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/external"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var productTracer = otel.Tracer("bitedash/internal/service/product")

type ProductService struct {
	queries        productQuerier
	fetchMenuItems menuItemsFetcher
}

func NewProductService(q productQuerier) *ProductService {
	return &ProductService{
		queries:        q,
		fetchMenuItems: external.FetchAllMenuItems,
	}
}

func (s *ProductService) SyncProducts(ctx context.Context) error {
	ctx, span := productTracer.Start(ctx, "ProductService.SyncProducts")
	defer span.End()
	items, err := s.fetchMenuItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range items {
		restaurant, err := s.queries.GetRestaurantByExternalID(ctx, int32(item.RestaurantID))
		if err != nil {
			return err
		}
		err = s.queries.UpsertProduct(ctx, db.UpsertProductParams{
			ID:              uuid.New(),
			ExternalID:      sql.NullInt32{Int32: int32(item.ItemID), Valid: true},
			RestaurantID:    restaurant.Restaurantid,
			Name:            item.ItemName,
			Description:     sql.NullString{String: item.ItemDescription, Valid: true},
			Price:           strconv.FormatFloat(item.ItemPrice, 'f', 2, 64),
			ImageUrl:        sql.NullString{String: item.ImageURL, Valid: item.ImageURL != ""},
			IsAvailable:     true,
			KitchenQuantity: 100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
