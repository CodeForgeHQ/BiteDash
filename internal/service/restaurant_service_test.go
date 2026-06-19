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

func TestRestaurantService_ListRestaurants_paginationDefaults(t *testing.T) {
	t.Parallel()

	var captured db.ListRestaurantsParams
	mock := &mockRestaurantQuerier{
		listRestaurants: func(ctx context.Context, arg db.ListRestaurantsParams) ([]db.ListRestaurantsRow, error) {
			captured = arg
			return []db.ListRestaurantsRow{}, nil
		},
	}

	svc := NewRestaurantService(mock)
	_, err := svc.ListRestaurants(context.Background(), "pizza", "italian", 0, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(20), captured.Limit)
	assert.Equal(t, int32(0), captured.Offset)
	assert.Equal(t, "pizza", captured.Column1)
	assert.Equal(t, "italian", captured.Column2)
}

func TestRestaurantService_ListRestaurants_capsLimit(t *testing.T) {
	t.Parallel()

	var captured db.ListRestaurantsParams
	mock := &mockRestaurantQuerier{
		listRestaurants: func(ctx context.Context, arg db.ListRestaurantsParams) ([]db.ListRestaurantsRow, error) {
			captured = arg
			return nil, nil
		},
	}

	svc := NewRestaurantService(mock)
	_, err := svc.ListRestaurants(context.Background(), "", "", 2, 100)
	require.NoError(t, err)
	assert.Equal(t, int32(20), captured.Limit)
	assert.Equal(t, int32(20), captured.Offset)
}

func TestRestaurantService_GetRestaurantDetails_success(t *testing.T) {
	t.Parallel()

	restaurantID := uuid.New()
	productID := uuid.New()

	mock := &mockRestaurantQuerier{
		getRestaurantWithProducts: func(ctx context.Context, id uuid.UUID) ([]db.GetRestaurantWithProductsRow, error) {
			assert.Equal(t, restaurantID, id)
			return []db.GetRestaurantWithProductsRow{
				{
					Restaurantid:   restaurantID,
					Restaurantname: "Pizza Place",
					Category:       "italian",
					Address:        "Main St",
					Parkinglot:     true,
					ProductID:      uuid.NullUUID{UUID: productID, Valid: true},
					ProductName:    sql.NullString{String: "Margherita", Valid: true},
					Description:    sql.NullString{String: "Classic", Valid: true},
					Price:          sql.NullString{String: "12.50", Valid: true},
					IsAvailable:    sql.NullBool{Bool: true, Valid: true},
				},
			}, nil
		},
	}

	svc := NewRestaurantService(mock)
	details, err := svc.GetRestaurantDetails(context.Background(), restaurantID)
	require.NoError(t, err)
	assert.Equal(t, restaurantID.String(), details.ID)
	assert.Equal(t, "Pizza Place", details.Name)
	require.Len(t, details.Products, 1)
	assert.Equal(t, productID.String(), details.Products[0].ID)
	assert.Equal(t, 12.5, details.Products[0].Price)
	assert.True(t, details.Products[0].Available)
}

func TestRestaurantService_GetRestaurantDetails_notFound(t *testing.T) {
	t.Parallel()

	mock := &mockRestaurantQuerier{
		getRestaurantWithProducts: func(ctx context.Context, id uuid.UUID) ([]db.GetRestaurantWithProductsRow, error) {
			return []db.GetRestaurantWithProductsRow{}, nil
		},
	}

	svc := NewRestaurantService(mock)
	_, err := svc.GetRestaurantDetails(context.Background(), uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestRestaurantService_GetRestaurantDetails_invalidPrice(t *testing.T) {
	t.Parallel()

	restaurantID := uuid.New()
	mock := &mockRestaurantQuerier{
		getRestaurantWithProducts: func(ctx context.Context, id uuid.UUID) ([]db.GetRestaurantWithProductsRow, error) {
			return []db.GetRestaurantWithProductsRow{
				{
					Restaurantid:   restaurantID,
					Restaurantname: "Cafe",
					Category:       "cafe",
					Address:        "Side St",
					ProductID:      uuid.NullUUID{UUID: uuid.New(), Valid: true},
					ProductName:    sql.NullString{String: "Tea", Valid: true},
					Price:          sql.NullString{String: "not-a-price", Valid: true},
					IsAvailable:    sql.NullBool{Bool: true, Valid: true},
				},
			}, nil
		},
	}

	svc := NewRestaurantService(mock)
	details, err := svc.GetRestaurantDetails(context.Background(), restaurantID)
	require.NoError(t, err)
	require.Len(t, details.Products, 1)
	assert.Equal(t, float64(0), details.Products[0].Price)
}

func TestRestaurantService_SyncRestaurants_success(t *testing.T) {
	t.Parallel()

	var upsertCount int
	mock := &mockRestaurantQuerier{
		upsertRestaurant: func(ctx context.Context, arg db.UpsertRestaurantParams) error {
			upsertCount++
			assert.Equal(t, int32(42), arg.Externalid)
			assert.Equal(t, "Test Bistro", arg.Restaurantname)
			return nil
		},
	}

	svc := NewRestaurantService(mock)
	svc.fetchRestaurants = func(ctx context.Context) ([]external.Restaurant, error) {
		return []external.Restaurant{
			{
				RestaurantID:   42,
				RestaurantName: "Test Bistro",
				Description:    "Nice food",
				Address:        "1 Test Ave",
				Type:           "bistro",
				ParkingLot:     true,
			},
		}, nil
	}

	err := svc.SyncRestaurants(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, upsertCount)
}

func TestRestaurantService_SyncRestaurants_fetchError(t *testing.T) {
	t.Parallel()

	fetchErr := errors.New("api down")
	svc := NewRestaurantService(&mockRestaurantQuerier{})
	svc.fetchRestaurants = func(ctx context.Context) ([]external.Restaurant, error) {
		return nil, fetchErr
	}

	err := svc.SyncRestaurants(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, fetchErr)
}

func TestRestaurantService_SyncRestaurants_upsertError(t *testing.T) {
	t.Parallel()

	upsertErr := errors.New("db write failed")
	mock := &mockRestaurantQuerier{
		upsertRestaurant: func(ctx context.Context, arg db.UpsertRestaurantParams) error {
			return upsertErr
		},
	}

	svc := NewRestaurantService(mock)
	svc.fetchRestaurants = func(ctx context.Context) ([]external.Restaurant, error) {
		return []external.Restaurant{{RestaurantID: 1, RestaurantName: "A"}}, nil
	}

	err := svc.SyncRestaurants(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, upsertErr)
}
