package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

func runListRestaurants(args []string) error {
	fs := flag.NewFlagSet("list-restaurants", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	page := fs.Int("page", 1, "page number")
	limit := fs.Int("limit", 10, "items per page")
	category := fs.String("category", "", "restaurant category")
	search := fs.String("search", "", "search query")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken("", *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.restaurantClient.ListRestaurants(ctx, &bitedashv1.ListRestaurantsRequest{
		Page:     int32(*page),
		Limit:    int32(*limit),
		Category: *category,
		Search:   *search,
	})
	if err != nil {
		return formatGRPCError("ListRestaurants", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))
	log.Printf("restaurants: page=%d limit=%d total=%d count=%d",
		resp.GetPage(),
		resp.GetLimit(),
		resp.GetTotal(),
		len(resp.GetRestaurants()),
	)

	for _, restaurant := range resp.GetRestaurants() {
		log.Printf("restaurant: id=%s name=%s category=%s address=%s products=%d",
			restaurant.GetId(),
			restaurant.GetName(),
			restaurant.GetCategory(),
			restaurant.GetAddress(),
			len(restaurant.GetProducts()),
		)
	}

	return nil
}

func runGetRestaurant(args []string) error {
	fs := flag.NewFlagSet("get-restaurant", flag.ExitOnError)

	addr := fs.String("addr", defaultGRPCAddr, "gRPC server address")
	restaurantID := fs.String("restaurant-id", "", "Restaurant ID")
	timeout := fs.Duration("timeout", 5*time.Second, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *restaurantID == "" {
		return fmt.Errorf("--restaurant-id is required")
	}

	clients, err := newGRPCClients(*addr)
	if err != nil {
		return fmt.Errorf("create grpc clients: %w", err)
	}
	defer clients.Close()

	ctx, cancel := contextWithToken("", *timeout)
	defer cancel()

	startedAt := time.Now()

	resp, err := clients.restaurantClient.GetRestaurantByID(ctx, &bitedashv1.GetRestaurantByIDRequest{
		RestaurantId: *restaurantID,
	})
	if err != nil {
		return formatGRPCError("GetRestaurantByID", err, time.Since(startedAt))
	}

	log.Printf("grpc call duration: %s", time.Since(startedAt))

	restaurant := resp.GetRestaurant()

	log.Printf("restaurant: id=%s name=%s category=%s address=%s parking_lot=%t products=%d",
		restaurant.GetId(),
		restaurant.GetName(),
		restaurant.GetCategory(),
		restaurant.GetAddress(),
		restaurant.GetParkingLot(),
		len(restaurant.GetProducts()),
	)

	for _, product := range restaurant.GetProducts() {
		log.Printf("  product: id=%s name=%s price=%.2f available=%t",
			product.GetId(),
			product.GetName(),
			product.GetPrice(),
			product.GetIsAvailable(),
		)
	}

	return nil
}
