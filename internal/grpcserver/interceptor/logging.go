package interceptor

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startedAt)

		attrs := []any{
			"method", info.FullMethod,
			"duration", duration.String(),
		}

		if requestID, ok := RequestIDFromContext(ctx); ok {
			attrs = append(attrs, "request_id", requestID)
		}

		if userID, ok := UserIDFromContext(ctx); ok {
			attrs = append(attrs, "user_id", userID.String())
		}

		if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
			attrs = append(attrs,
				"trace_id", spanCtx.TraceID().String(),
				"span_id", spanCtx.SpanID().String(),
			)
		}

		if err != nil {
			st, _ := status.FromError(err)

			attrs = append(attrs,
				"code", st.Code().String(),
				"message", st.Message(),
			)

			logger.Warn("grpc request failed", attrs...)

			return nil, err
		}

		logger.Info("grpc request completed", attrs...)

		return resp, nil
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

		ctx := stream.Context()

		attrs := []any{
			"method", info.FullMethod,
			"duration", duration.String(),
		}

		if requestID, ok := RequestIDFromContext(ctx); ok {
			attrs = append(attrs, "request_id", requestID)
		}

		if userID, ok := UserIDFromContext(ctx); ok {
			attrs = append(attrs, "user_id", userID.String())
		}

		if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
			attrs = append(attrs,
				"trace_id", spanCtx.TraceID().String(),
				"span_id", spanCtx.SpanID().String(),
			)
		}

		if err != nil {
			st, _ := status.FromError(err)

			attrs = append(attrs,
				"code", st.Code().String(),
				"message", st.Message(),
			)

			logger.Warn("grpc stream failed", attrs...)

			return err
		}

		logger.Info("grpc stream completed", attrs...)

		return nil
	}
}
