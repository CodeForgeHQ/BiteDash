package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/dto"
	"bitedash/internal/external"

	"github.com/google/uuid"
)

type RestaurantService struct {
	queries          restaurantQuerier
	fetchRestaurants restaurantFetcher
}

func NewRestaurantService(q restaurantQuerier) *RestaurantService {
	return &RestaurantService{
		queries:          q,
		fetchRestaurants: external.FetchRestaurants,
	}
}

func (s *RestaurantService) SyncRestaurants(ctx context.Context) error {

	apiRestaurants, err := s.fetchRestaurants(ctx)
	if err != nil {
		return err
	}

	for _, r := range apiRestaurants {

		err = s.queries.UpsertRestaurant(ctx, db.UpsertRestaurantParams{
			Restaurantid:   uuid.New(),
			Externalid:     int32(r.RestaurantID),
			Restaurantname: r.RestaurantName,
			Description:    sql.NullString{String: r.Description, Valid: true},
			Category:       r.Type,
			Address:        r.Address,
			Parkinglot:     r.ParkingLot,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *RestaurantService) ListRestaurants(ctx context.Context, search, category string, page, limit int32) (*dto.ListRestaurantsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	rows, err := s.queries.ListRestaurants(ctx, db.ListRestaurantsParams{
		Column1: search,
		Column2: category,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountRestaurants(ctx, db.CountRestaurantsParams{
		Column1: search,
		Column2: category,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.ListRestaurantsResponse{
		Restaurants: make([]dto.RestaurantDetails, 0, len(rows)),
		Total:       total,
		Page:        page,
		Limit:       limit,
	}
	for _, row := range rows {
		response.Restaurants = append(response.Restaurants, dto.RestaurantDetails{
			ID:       row.Restaurantid.String(),
			Name:     row.Restaurantname,
			Category: row.Category,
			Address:  row.Address,
			Parking:  row.Parkinglot,
			Products: make([]dto.Product, 0),
		})
	}

	return response, nil
}

func (s *RestaurantService) GetRestaurantDetails(ctx context.Context, restaurantID uuid.UUID) (*dto.RestaurantDetails, error) {
	rows, err := s.queries.GetRestaurantWithProducts(ctx, restaurantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRestaurantNotFound
		}
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrRestaurantNotFound
	}

	restaurant := &dto.RestaurantDetails{
		ID:       rows[0].Restaurantid.String(),
		Name:     rows[0].Restaurantname,
		Category: rows[0].Category,
		Address:  rows[0].Address,
		Parking:  rows[0].Parkinglot,
		Products: make([]dto.Product, 0),
	}

	for _, row := range rows {
		if row.ProductID.Valid {
			restaurant.Products = append(restaurant.Products, dto.Product{
				ID:          row.ProductID.UUID.String(),
				Name:        row.ProductName.String,
				Description: row.Description.String,
				Price:       parsePrice(row.Price),
				Available:   row.IsAvailable.Bool,
			})
		}
	}

	return restaurant, nil
}

// parsePrice safely converts sql.NullString to float64
func parsePrice(ns sql.NullString) float64 {
	if !ns.Valid {
		return 0
	}
	f, err := strconv.ParseFloat(ns.String, 64)
	if err != nil {
		return 0
	}
	return f
}
