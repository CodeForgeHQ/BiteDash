package service

import (
	"context"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/external"

	"github.com/google/uuid"
)

type authQuerier interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type restaurantQuerier interface {
	UpsertRestaurant(ctx context.Context, arg db.UpsertRestaurantParams) error
	ListRestaurants(ctx context.Context, arg db.ListRestaurantsParams) ([]db.ListRestaurantsRow, error)
	GetRestaurantWithProducts(ctx context.Context, restaurantid uuid.UUID) ([]db.GetRestaurantWithProductsRow, error)
}

type restaurantFetcher func(ctx context.Context) ([]external.Restaurant, error)

type productQuerier interface {
	GetRestaurantByExternalID(ctx context.Context, externalid int32) (db.Restaurant, error)
	UpsertProduct(ctx context.Context, arg db.UpsertProductParams) error
}

type menuItemsFetcher func(ctx context.Context) ([]external.MenuItem, error)
