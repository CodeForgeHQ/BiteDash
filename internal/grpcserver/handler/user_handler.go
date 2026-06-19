package handler

import (
	"context"
	"errors"
	"log/slog"

	"bitedash/internal/grpcserver/interceptor"
	"bitedash/internal/grpcserver/mapper"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
	"bitedash/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	bitedashv1.UnimplementedUserServiceServer

	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

func (h *UserHandler) GetUserByID(
	ctx context.Context,
	req *bitedashv1.GetUserByIDRequest,
) (*bitedashv1.GetUserByIDResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, mapUserServiceError(err)
	}

	return &bitedashv1.GetUserByIDResponse{
		User: mapper.UserToProto(user),
	}, nil
}

func mapUserServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	default:
		slog.Error("unexpected user service error", "error", err)
		return status.Error(codes.Internal, "internal server error")
	}
}

/*func simulateSlowWork(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}*/

func (h *UserHandler) GetMe(
	ctx context.Context,
	req *bitedashv1.GetMeRequest,
) (*bitedashv1.GetMeResponse, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id is missing in context")
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, mapUserServiceError(err)
	}

	return &bitedashv1.GetMeResponse{
		User: mapper.UserToProto(user),
	}, nil
}
