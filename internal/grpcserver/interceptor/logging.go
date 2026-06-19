package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)

		st := status.Code(err)

		slog.Info(
			"grpc request completed",
			"method", info.FullMethod,
			"code", st.String(),
			"duration", time.Since(startedAt).String(),
		)

		return resp, err
	}
}

func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startedAt := time.Now()

		err := handler(srv, stream)

		duration := time.Since(startedAt)

		if err != nil {
			st, _ := status.FromError(err)

			logger.Warn("grpc stream failed",
				"method", info.FullMethod,
				"code", st.Code().String(),
				"message", st.Message(),
				"duration", duration.String(),
			)

			return err
		}

		logger.Info("grpc stream completed",
			"method", info.FullMethod,
			"duration", duration.String(),
		)

		return nil
	}
}
