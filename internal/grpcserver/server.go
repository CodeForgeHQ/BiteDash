package grpcserver

import (
	"log/slog"

	"bitedash/internal/grpcserver/handler"
	"bitedash/internal/grpcserver/interceptor"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"

	"google.golang.org/grpc"
)

type Deps struct {
	UserHandler       *handler.UserHandler
	OrderHandler      *handler.OrderHandler
	RestaurantHandler *handler.RestaurantHandler
	CartHandler       *handler.CartHandler
}

type Server struct {
	server *grpc.Server
}

func NewServer(deps Deps) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnaryInterceptor(slog.Default()),
			interceptor.RequestIDUnaryInterceptor(),
			interceptor.AuthUnaryInterceptor(),
			interceptor.LoggingUnaryInterceptor(slog.Default()),
		),
		grpc.ChainStreamInterceptor(
			interceptor.RecoveryStreamInterceptor(slog.Default()),
			interceptor.RequestIDStreamInterceptor(),
			interceptor.AuthStreamInterceptor(),
			interceptor.LoggingStreamInterceptor(slog.Default()),
		),
	)

	bitedashv1.RegisterUserServiceServer(grpcServer, deps.UserHandler)
	bitedashv1.RegisterOrderServiceServer(grpcServer, deps.OrderHandler)
	bitedashv1.RegisterRestaurantServiceServer(grpcServer, deps.RestaurantHandler)
	bitedashv1.RegisterCartServiceServer(grpcServer, deps.CartHandler)
	return &Server{
		server: grpcServer,
	}
}

func (s *Server) Server() *grpc.Server {
	return s.server
}
