package grpcserver

import (
	"bitedash/internal/grpcserver/handler"
	"bitedash/internal/grpcserver/interceptor"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
	"log/slog"

	"google.golang.org/grpc"
)

type Deps struct {
	UserHandler  *handler.UserHandler
	OrderHandler *handler.OrderHandler
}

type Server struct {
	server *grpc.Server
}

func NewServer(deps Deps) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnaryInterceptor(),
			interceptor.LoggingUnaryInterceptor(),
			interceptor.AuthUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.RecoveryStreamInterceptor(slog.Default()),
			interceptor.LoggingStreamInterceptor(slog.Default()),
			interceptor.AuthStreamInterceptor(),
		),
	)

	bitedashv1.RegisterUserServiceServer(grpcServer, deps.UserHandler)
	bitedashv1.RegisterOrderServiceServer(grpcServer, deps.OrderHandler)
	return &Server{
		server: grpcServer,
	}
}

func (s *Server) Server() *grpc.Server {
	return s.server
}
