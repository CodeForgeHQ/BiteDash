package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

func runCheckout(args []string) error {
	fs := flag.NewFlagSet("checkout", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.orderClient.Checkout(ctx, &bitedashv1.CheckoutRequest{})
	if err != nil {
		return formatGRPCError("Checkout", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	log.Printf("checkout completed: order_id=%s status=%s total_amount=%.2f items_count=%d",
		resp.GetOrderId(),
		resp.GetStatus(),
		resp.GetTotalAmount(),
		resp.GetItemsCount(),
	)

	return nil
}

func runListOrders(args []string) error {
	fs := flag.NewFlagSet("list-orders", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.orderClient.ListMyOrders(ctx, &bitedashv1.ListMyOrdersRequest{})
	if err != nil {
		return formatGRPCError("ListMyOrders", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))

	if len(resp.GetOrders()) == 0 {
		log.Println("orders: empty")
		return nil
	}

	for _, order := range resp.GetOrders() {
		printOrder(order)
	}

	return nil
}

func runGetOrder(args []string) error {
	fs := flag.NewFlagSet("get-order", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	orderID := fs.String("order-id", "", "Order ID")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	if *orderID == "" {
		return fmt.Errorf("--order-id is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.orderClient.GetOrderByID(ctx, &bitedashv1.GetOrderByIDRequest{
		OrderId: *orderID,
	})
	if err != nil {
		return formatGRPCError("GetOrderByID", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	printOrder(resp.GetOrder())

	return nil
}

func runWatchOrder(args []string) error {
	fs := flag.NewFlagSet("watch-order", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	orderID := fs.String("order-id", "", "Order ID")
	timeout := fs.Duration("timeout", 10*time.Second, "stream timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	if *orderID == "" {
		return fmt.Errorf("--order-id is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	stream, err := clients.orderClient.WatchOrderStatus(ctx, &bitedashv1.WatchOrderStatusRequest{
		OrderId: *orderID,
	})
	if err != nil {
		return formatGRPCError("WatchOrderStatus", err, 0)
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("stream completed")
				return nil
			}

			return formatGRPCError("WatchOrderStatus stream", err, 0)
		}

		log.Printf("order event: order_id=%s status=%s message=%s occurred_at=%s",
			event.GetOrderId(),
			event.GetStatus(),
			event.GetMessage(),
			event.GetOccurredAt(),
		)
	}
}

func printOrder(order *bitedashv1.Order) {
	if order == nil {
		log.Println("order: nil")
		return
	}

	log.Printf("order: id=%s status=%s total=%.2f created_at=%s items=%d",
		order.GetId(),
		order.GetStatus(),
		order.GetTotalAmount(),
		order.GetCreatedAt(),
		len(order.GetItems()),
	)

	for _, item := range order.GetItems() {
		log.Printf("  item: product=%s name=%s unit_price=%.2f qty=%d line_total=%.2f",
			item.GetProductId(),
			item.GetProductName(),
			item.GetUnitPrice(),
			item.GetQuantity(),
			item.GetLineTotal(),
		)
	}
}
