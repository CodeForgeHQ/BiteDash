package handler

import (
	"context"
	"errors"

	"bitedash/internal/grpcserver/interceptor"
	"bitedash/internal/grpcserver/mapper"
	"bitedash/internal/grpcserver/validation"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
	"bitedash/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartHandler struct {
	bitedashv1.UnimplementedCartServiceServer

	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) AddItem(
	ctx context.Context,
	req *bitedashv1.AddCartItemRequest,
) (*bitedashv1.CartResponse, error) {
	userID, err := cartUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	request, err := validation.ValidateAddCartItemRequest(req)
	if err != nil {
		return nil, err
	}

	err = h.cartService.AddItem(ctx, userID, request)
	if err != nil {
		return nil, mapCartServiceError(err)
	}

	return h.getCart(ctx, userID)
}

func (h *CartHandler) GetCart(
	ctx context.Context,
	req *bitedashv1.GetCartRequest,
) (*bitedashv1.CartResponse, error) {
	userID, err := cartUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return h.getCart(ctx, userID)
}

func (h *CartHandler) RemoveItem(
	ctx context.Context,
	req *bitedashv1.RemoveCartItemRequest,
) (*bitedashv1.CartResponse, error) {
	userID, err := cartUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	productID, err := parseCartProductID(req.GetProductId())
	if err != nil {
		return nil, err
	}

	if err := h.cartService.RemoveCartItem(ctx, userID, productID); err != nil {
		return nil, mapCartServiceError(err)
	}

	return h.getCart(ctx, userID)
}

func (h *CartHandler) ClearCart(
	ctx context.Context,
	req *bitedashv1.ClearCartRequest,
) (*bitedashv1.ClearCartResponse, error) {
	userID, err := cartUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.cartService.ClearCart(ctx, userID); err != nil {
		return nil, mapCartServiceError(err)
	}

	return &bitedashv1.ClearCartResponse{Success: true}, nil
}

func (h *CartHandler) getCart(ctx context.Context, userID uuid.UUID) (*bitedashv1.CartResponse, error) {
	cart, err := h.cartService.GetCart(ctx, userID)
	if err != nil {
		return nil, mapCartServiceError(err)
	}

	return mapper.CartResponseToProto(cart), nil
}

func cartUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	return userID, nil
}

func parseCartProductID(value string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	productID, err := uuid.Parse(value)
	if err != nil || productID == uuid.Nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid product_id")
	}

	return productID, nil
}

func mapCartServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrProductNotFound):
		return status.Error(codes.NotFound, "product not found")
	case errors.Is(err, service.ErrCartNotFound):
		return status.Error(codes.NotFound, "active cart not found")
	case errors.Is(err, service.ErrCartItemNotFound):
		return status.Error(codes.NotFound, "cart item not found")
	case errors.Is(err, service.ErrProductNotAvailable):
		return status.Error(codes.FailedPrecondition, "product is not available")
	case errors.Is(err, service.ErrInsufficientStock):
		return status.Error(codes.FailedPrecondition, "insufficient product stock")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
