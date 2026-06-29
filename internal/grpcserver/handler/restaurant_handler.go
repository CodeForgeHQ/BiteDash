package handler

import (
	"context"
	"errors"

	"bitedash/internal/grpcserver/mapper"
	"bitedash/internal/grpcserver/validation"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
	"bitedash/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RestaurantHandler struct {
	bitedashv1.UnimplementedRestaurantServiceServer

	restaurantService *service.RestaurantService
}

func NewRestaurantHandler(restaurantService *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

func (h *RestaurantHandler) ListRestaurants(
	ctx context.Context,
	req *bitedashv1.ListRestaurantsRequest,
) (*bitedashv1.ListRestaurantsResponse, error) {
	if err := validation.ValidateListRestaurantsRequest(req); err != nil {
		return nil, err
	}

	response, err := h.restaurantService.ListRestaurants(
		ctx,
		req.GetSearch(),
		req.GetCategory(),
		req.GetPage(),
		req.GetLimit(),
	)
	if err != nil {
		return nil, mapRestaurantServiceError(err)
	}

	return mapper.ListRestaurantsResponseToProto(*response), nil
}

func (h *RestaurantHandler) GetRestaurantByID(
	ctx context.Context,
	req *bitedashv1.GetRestaurantByIDRequest,
) (*bitedashv1.GetRestaurantByIDResponse, error) {
	if req.GetRestaurantId() == "" {
		return nil, status.Error(codes.InvalidArgument, "restaurant_id is required")
	}

	restaurantID, err := uuid.Parse(req.GetRestaurantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid restaurant_id")
	}

	restaurant, err := h.restaurantService.GetRestaurantDetails(ctx, restaurantID)
	if err != nil {
		return nil, mapRestaurantServiceError(err)
	}

	return &bitedashv1.GetRestaurantByIDResponse{
		Restaurant: mapper.RestaurantToProto(*restaurant),
	}, nil
}

func mapRestaurantServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrRestaurantNotFound):
		return status.Error(codes.NotFound, "restaurant not found")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
