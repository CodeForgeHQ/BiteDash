package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"

	"google.golang.org/grpc/status"
)

func runGetMe(args []string) error {
	fs := flag.NewFlagSet("get-me", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	timeout := fs.Duration("timeout", 2*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.userClient.GetMe(ctx, &bitedashv1.GetMeRequest{})
	if err != nil {
		return formatGRPCError("GetMe", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	log.Printf("me: id=%s email=%s created_at=%s",
		resp.GetUser().GetId(),
		resp.GetUser().GetEmail(),
		resp.GetUser().GetCreatedAt(),
	)

	return nil
}

func runGetUser(args []string) error {
	fs := flag.NewFlagSet("get-user", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	userID := fs.String("user-id", "", "User ID")
	timeout := fs.Duration("timeout", 2*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *userID == "" {
		return fmt.Errorf("--user-id is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken("", *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.userClient.GetUserByID(ctx, &bitedashv1.GetUserByIDRequest{
		UserId: *userID,
	})
	if err != nil {
		return formatGRPCError("GetUserByID", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	log.Printf("user: id=%s email=%s created_at=%s",
		resp.GetUser().GetId(),
		resp.GetUser().GetEmail(),
		resp.GetUser().GetCreatedAt(),
	)

	return nil
}

func formatGRPCError(method string, err error, duration time.Duration) error {
	st, ok := status.FromError(err)
	if ok {
		return fmt.Errorf("%s failed: code=%s message=%s duration=%s",
			method,
			st.Code(),
			st.Message(),
			duration,
		)
	}

	return fmt.Errorf("%s failed: %w", method, err)
}
