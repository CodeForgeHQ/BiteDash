package interceptor

import (
	"context"
	"strings"

	authpkg "bitedash/internal/pkg/auth"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/bitedash.v1.UserService/GetUserByID": {},
}

func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		userID, err := userIDFromMetadata(ctx)
		if err != nil {
			return nil, err
		}

		ctx = ContextWithUserID(ctx, userID)

		return handler(ctx, req)
	}
}

func isPublicMethod(fullMethod string) bool {
	switch fullMethod {
	case "/bitedash.v1.UserService/GetUserByID",
		"/bitedash.v1.RestaurantService/ListRestaurants",
		"/bitedash.v1.RestaurantService/GetRestaurantByID":
		return true
	default:
		return false
	}
}

func userIDFromMetadata(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "metadata is missing")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return uuid.Nil, status.Error(codes.Unauthenticated, "authorization metadata is missing")
	}

	authHeader := values[0]

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
	if tokenString == "" {
		return uuid.Nil, status.Error(codes.Unauthenticated, "token is empty")
	}

	userID, err := authpkg.ParseUserIDFromToken(tokenString)
	if err != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return userID, nil
}

func AuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if isPublicMethod(info.FullMethod) {
			return handler(srv, stream)
		}

		userID, err := userIDFromMetadata(stream.Context())
		if err != nil {
			return err
		}

		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ContextWithUserID(stream.Context(), userID),
		}

		return handler(srv, wrappedStream)
	}
}
