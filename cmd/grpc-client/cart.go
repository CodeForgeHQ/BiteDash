package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

func runAddCartItem(args []string) error {
	fs := flag.NewFlagSet("add-cart-item", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	productID := fs.String("product-id", "", "Product ID")
	quantity := fs.Int("quantity", 1, "Product quantity")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	if *productID == "" {
		return fmt.Errorf("--product-id is required")
	}

	if *quantity <= 0 {
		return fmt.Errorf("--quantity must be greater than 0")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.cartClient.AddItem(ctx, &bitedashv1.AddCartItemRequest{
		ProductId: *productID,
		Quantity:  int32(*quantity),
	})
	if err != nil {
		return formatGRPCError("AddItem", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	printCart(resp)

	return nil
}

func runGetCart(args []string) error {
	fs := flag.NewFlagSet("get-cart", flag.ExitOnError)

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

	resp, err := clients.cartClient.GetCart(ctx, &bitedashv1.GetCartRequest{})
	if err != nil {
		return formatGRPCError("GetCart", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	printCart(resp)

	return nil
}

func runRemoveCartItem(args []string) error {
	fs := flag.NewFlagSet("remove-cart-item", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	token := fs.String("token", "", "JWT access token")
	productID := fs.String("product-id", "", "Product ID")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *token == "" {
		return fmt.Errorf("--token is required")
	}

	if *productID == "" {
		return fmt.Errorf("--product-id is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken(*token, *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.cartClient.RemoveItem(ctx, &bitedashv1.RemoveCartItemRequest{
		ProductId: *productID,
	})
	if err != nil {
		return formatGRPCError("RemoveItem", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	printCart(resp)

	return nil
}

func runClearCart(args []string) error {
	fs := flag.NewFlagSet("clear-cart", flag.ExitOnError)

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

	resp, err := clients.cartClient.ClearCart(ctx, &bitedashv1.ClearCartRequest{})
	if err != nil {
		return formatGRPCError("ClearCart", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	log.Printf("cart cleared: success=%t", resp.GetSuccess())

	return nil
}

func printCart(cart *bitedashv1.CartResponse) {
	if cart == nil {
		log.Println("cart: nil")
		return
	}

	log.Printf("cart: id=%s total=%.2f items=%d",
		cart.GetCartId(),
		cart.GetTotalAmount(),
		len(cart.GetItems()),
	)

	for _, item := range cart.GetItems() {
		log.Printf("  item: product=%s name=%s unit_price=%.2f qty=%d line_total=%.2f",
			item.GetProductId(),
			item.GetProductName(),
			item.GetUnitPrice(),
			item.GetQuantity(),
			item.GetLineTotal(),
		)
	}
}
