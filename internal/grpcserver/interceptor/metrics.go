package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"bitedash/internal/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startedAt)

		code := codes.OK
		if err != nil {
			st, _ := status.FromError(err)
			code = st.Code()
		}
		metrics.ObserveGRPCRequest(info.FullMethod, code, duration)

		return resp, err
	}
}

func MetricsStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startedAt := time.Now()

		err := handler(srv, stream)

		duration := time.Since(startedAt)

		code := codes.OK
		if err != nil {
			st, _ := status.FromError(err)
			code = st.Code()
		}
		metrics.ObserveGRPCRequest(info.FullMethod, code, duration)

		return err
	}
}
