package interceptor

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIDMetadataKey = "x-request-id"

func RequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID := requestIDFromMetadata(ctx)
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx = ContextWithRequestID(ctx, requestID)

		_ = grpc.SetHeader(ctx, metadata.Pairs(requestIDMetadataKey, requestID))

		return handler(ctx, req)
	}
}

func RequestIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		requestID := requestIDFromMetadata(stream.Context())
		if requestID == "" {
			requestID = newRequestID()
		}

		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ContextWithRequestID(stream.Context(), requestID),
		}

		_ = wrappedStream.SendHeader(metadata.Pairs(requestIDMetadataKey, requestID))

		return handler(srv, wrappedStream)
	}
}

func requestIDFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(requestIDMetadataKey)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func newRequestID() string {
	var b [8]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(b[:])
}
