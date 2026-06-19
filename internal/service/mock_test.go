package service

import (
	"context"
	"errors"

	db "bitedash/internal/db/sqlc"

	"github.com/google/uuid"
)

var errMockNotConfigured = errors.New("mock: handler not configured")

type mockRestaurantQuerier struct {
	upsertRestaurant          func(ctx context.Context, arg db.UpsertRestaurantParams) error
	listRestaurants           func(ctx context.Context, arg db.ListRestaurantsParams) ([]db.ListRestaurantsRow, error)
	getRestaurantWithProducts func(ctx context.Context, restaurantid uuid.UUID) ([]db.GetRestaurantWithProductsRow, error)
}

func (m *mockRestaurantQuerier) UpsertRestaurant(ctx context.Context, arg db.UpsertRestaurantParams) error {
	if m.upsertRestaurant != nil {
		return m.upsertRestaurant(ctx, arg)
	}
	return errMockNotConfigured
}

func (m *mockRestaurantQuerier) ListRestaurants(ctx context.Context, arg db.ListRestaurantsParams) ([]db.ListRestaurantsRow, error) {
	if m.listRestaurants != nil {
		return m.listRestaurants(ctx, arg)
	}
	return nil, errMockNotConfigured
}

func (m *mockRestaurantQuerier) GetRestaurantWithProducts(ctx context.Context, restaurantid uuid.UUID) ([]db.GetRestaurantWithProductsRow, error) {
	if m.getRestaurantWithProducts != nil {
		return m.getRestaurantWithProducts(ctx, restaurantid)
	}
	return nil, errMockNotConfigured
}

type mockProductQuerier struct {
	getRestaurantByExternalID func(ctx context.Context, externalid int32) (db.Restaurant, error)
	upsertProduct             func(ctx context.Context, arg db.UpsertProductParams) error
}

func (m *mockProductQuerier) GetRestaurantByExternalID(ctx context.Context, externalid int32) (db.Restaurant, error) {
	if m.getRestaurantByExternalID != nil {
		return m.getRestaurantByExternalID(ctx, externalid)
	}
	return db.Restaurant{}, errMockNotConfigured
}

func (m *mockProductQuerier) UpsertProduct(ctx context.Context, arg db.UpsertProductParams) error {
	if m.upsertProduct != nil {
		return m.upsertProduct(ctx, arg)
	}
	return errMockNotConfigured
}
