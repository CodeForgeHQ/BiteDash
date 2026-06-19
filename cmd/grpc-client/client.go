package main

import (
	"context"
	"time"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const defaultGRPCAddr = "127.0.0.1:9090"

type grpcClients struct {
	conn        *grpc.ClientConn
	userClient  bitedashv1.UserServiceClient
	orderClient bitedashv1.OrderServiceClient
}

func newGRPCClients(addr string) (*grpcClients, error) {
	if addr == "" {
		addr = defaultGRPCAddr
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		conn:        conn,
		userClient:  bitedashv1.NewUserServiceClient(conn),
		orderClient: bitedashv1.NewOrderServiceClient(conn),
	}, nil
}

func (c *grpcClients) Close() error {
	return c.conn.Close()
}

func contextWithToken(token string, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	if token == "" {
		return ctx, cancel
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"authorization",
		"Bearer "+token,
	)

	return ctx, cancel
}
