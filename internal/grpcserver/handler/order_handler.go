package handler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"bitedash/internal/grpcserver/interceptor"
	"bitedash/internal/grpcserver/mapper"
	"bitedash/internal/grpcserver/validation"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
	"bitedash/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	bitedashv1.UnimplementedOrderServiceServer

	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Checkout(
	ctx context.Context,
	req *bitedashv1.CheckoutRequest,
) (*bitedashv1.CheckoutResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	order, err := h.orderService.MakeOrder(ctx, userID)
	if err != nil {
		return nil, mapOrderServiceError(err)
	}

	return &bitedashv1.CheckoutResponse{
		OrderId:     order.OrderID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		ItemsCount:  order.ItemsCount,
	}, nil
}

func (h *OrderHandler) ListMyOrders(
	ctx context.Context,
	req *bitedashv1.ListMyOrdersRequest,
) (*bitedashv1.ListMyOrdersResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	orders, err := h.orderService.ListMyOrders(ctx, userID)
	if err != nil {
		return nil, mapOrderServiceError(err)
	}

	result := make([]*bitedashv1.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, mapper.OrderToProto(order))
	}

	return &bitedashv1.ListMyOrdersResponse{
		Orders: result,
	}, nil
}

func (h *OrderHandler) GetOrderByID(
	ctx context.Context,
	req *bitedashv1.GetOrderByIDRequest,
) (*bitedashv1.GetOrderByIDResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	orderID, err := validation.ValidateGetOrderByIDRequest(req)
	if err != nil {
		return nil, err
	}

	order, err := h.orderService.GetOrderByID(ctx, userID, orderID)
	if err != nil {
		return nil, mapOrderServiceError(err)
	}

	return &bitedashv1.GetOrderByIDResponse{
		Order: mapper.OrderToProto(*order),
	}, nil
}

func (h *OrderHandler) WatchOrderStatus(
	req *bitedashv1.WatchOrderStatusRequest,
	stream bitedashv1.OrderService_WatchOrderStatusServer,
) error {
	ctx := stream.Context()

	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	if req.GetOrderId() == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid order_id")
	}

	// Проверяем, что заказ существует и принадлежит текущему пользователю.
	// Это важно: нельзя позволять смотреть события чужого заказа.
	order, err := h.orderService.GetOrderByID(ctx, userID, orderID)
	if err != nil {
		return mapOrderServiceError(err)
	}

	events := []struct {
		status  string
		message string
	}{
		{
			status:  order.Status,
			message: "Order status loaded",
		},
		{
			status:  "preparing",
			message: "Restaurant started preparing your order",
		},
		{
			status:  "delivering",
			message: "Courier is delivering your order",
		},
		{
			status:  "completed",
			message: "Order completed",
		},
	}

	for _, event := range events {
		select {
		case <-ctx.Done():
			slog.Warn("order status stream cancelled",
				"order_id", orderID.String(),
				"error", ctx.Err(),
			)
			return mapOrderServiceError(ctx.Err())

		case <-time.After(1 * time.Second):
			err := stream.Send(&bitedashv1.OrderStatusEvent{
				OrderId:    order.OrderID,
				Status:     event.status,
				Message:    event.message,
				OccurredAt: time.Now().Format(time.RFC3339),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mapOrderServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrCartNotFound):
		return status.Error(codes.FailedPrecondition, "active cart not found")
	case errors.Is(err, service.ErrCartEmpty):
		return status.Error(codes.FailedPrecondition, "cart is empty")
	case errors.Is(err, service.ErrOrderNotFound):
		return status.Error(codes.NotFound, "order not found")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
