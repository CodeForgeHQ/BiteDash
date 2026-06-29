package validation

import (
	"bitedash/internal/dto"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateGetOrderByIDRequest(req *bitedashv1.GetOrderByIDRequest) (uuid.UUID, error) {
	if req == nil || req.GetOrderId() == "" {
		return uuid.Nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil || orderID == uuid.Nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid order_id")
	}

	return orderID, nil
}

func ValidateAddCartItemRequest(req *bitedashv1.AddCartItemRequest) (dto.AddCartItemRequest, error) {
	if req == nil || req.GetProductId() == "" {
		return dto.AddCartItemRequest{}, status.Error(codes.InvalidArgument, "product_id is required")
	}

	productID, err := uuid.Parse(req.GetProductId())
	if err != nil || productID == uuid.Nil {
		return dto.AddCartItemRequest{}, status.Error(codes.InvalidArgument, "invalid product_id")
	}
	if req.GetQuantity() <= 0 {
		return dto.AddCartItemRequest{}, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	return dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  req.GetQuantity(),
	}, nil
}

func ValidateListRestaurantsRequest(req *bitedashv1.ListRestaurantsRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if req.GetPage() < 0 {
		return status.Error(codes.InvalidArgument, "page must not be negative")
	}
	if req.GetLimit() < 0 {
		return status.Error(codes.InvalidArgument, "limit must not be negative")
	}

	return nil
}
