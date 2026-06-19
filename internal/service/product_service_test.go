package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/external"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductService_SyncProducts_success(t *testing.T) {
	t.Parallel()

	restaurantID := uuid.New()
	var upserted db.UpsertProductParams

	mock := &mockProductQuerier{
		getRestaurantByExternalID: func(ctx context.Context, externalid int32) (db.Restaurant, error) {
			assert.Equal(t, int32(7), externalid)
			return db.Restaurant{Restaurantid: restaurantID, Externalid: 7}, nil
		},
		upsertProduct: func(ctx context.Context, arg db.UpsertProductParams) error {
			upserted = arg
			return nil
		},
	}

	svc := NewProductService(mock)
	svc.fetchMenuItems = func(context.Context) ([]external.MenuItem, error) {
		return []external.MenuItem{
			{
				ItemID:          101,
				ItemName:        "Burger",
				ItemDescription: "Tasty",
				ItemPrice:       9.99,
				RestaurantID:    7,
				ImageURL:        "https://example.com/burger.png",
			},
		}, nil
	}

	err := svc.SyncProducts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, restaurantID, upserted.RestaurantID)
	assert.Equal(t, "Burger", upserted.Name)
	assert.Equal(t, int32(101), upserted.ExternalID.Int32)
	assert.True(t, upserted.ExternalID.Valid)
	assert.Equal(t, "9.99", upserted.Price)
	assert.True(t, upserted.ImageUrl.Valid)
	assert.Equal(t, "https://example.com/burger.png", upserted.ImageUrl.String)
}

func TestProductService_SyncProducts_skipsEmptyImageURL(t *testing.T) {
	t.Parallel()

	mock := &mockProductQuerier{
		getRestaurantByExternalID: func(ctx context.Context, externalid int32) (db.Restaurant, error) {
			return db.Restaurant{Restaurantid: uuid.New()}, nil
		},
		upsertProduct: func(ctx context.Context, arg db.UpsertProductParams) error {
			assert.False(t, arg.ImageUrl.Valid)
			return nil
		},
	}

	svc := NewProductService(mock)
	svc.fetchMenuItems = func(context.Context) ([]external.MenuItem, error) {
		return []external.MenuItem{{ItemID: 1, ItemName: "Soup", RestaurantID: 1}}, nil
	}

	require.NoError(t, svc.SyncProducts(context.Background()))
}

func TestProductService_SyncProducts_fetchError(t *testing.T) {
	t.Parallel()

	fetchErr := errors.New("menu api unavailable")
	svc := NewProductService(&mockProductQuerier{})
	svc.fetchMenuItems = func(context.Context) ([]external.MenuItem, error) {
		return nil, fetchErr
	}

	err := svc.SyncProducts(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, fetchErr)
}

func TestProductService_SyncProducts_restaurantNotFound(t *testing.T) {
	t.Parallel()

	dbErr := sql.ErrNoRows
	mock := &mockProductQuerier{
		getRestaurantByExternalID: func(ctx context.Context, externalid int32) (db.Restaurant, error) {
			return db.Restaurant{}, dbErr
		},
	}

	svc := NewProductService(mock)
	svc.fetchMenuItems = func(context.Context) ([]external.MenuItem, error) {
		return []external.MenuItem{{ItemID: 1, RestaurantID: 99}}, nil
	}

	err := svc.SyncProducts(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
}

func TestProductService_SyncProducts_upsertError(t *testing.T) {
	t.Parallel()

	upsertErr := errors.New("upsert failed")
	mock := &mockProductQuerier{
		getRestaurantByExternalID: func(ctx context.Context, externalid int32) (db.Restaurant, error) {
			return db.Restaurant{Restaurantid: uuid.New()}, nil
		},
		upsertProduct: func(ctx context.Context, arg db.UpsertProductParams) error {
			return upsertErr
		},
	}

	svc := NewProductService(mock)
	svc.fetchMenuItems = func(context.Context) ([]external.MenuItem, error) {
		return []external.MenuItem{{ItemID: 1, RestaurantID: 1}}, nil
	}

	err := svc.SyncProducts(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, upsertErr)
}
