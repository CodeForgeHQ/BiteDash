package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	var err error

	switch command {
	case "get-me":
		err = runGetMe(os.Args[2:])
	case "get-user":
		err = runGetUser(os.Args[2:])
	case "checkout":
		err = runCheckout(os.Args[2:])
	case "list-orders":
		err = runListOrders(os.Args[2:])
	case "get-order":
		err = runGetOrder(os.Args[2:])
	case "watch-order":
		err = runWatchOrder(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Print(`BiteDash gRPC client

Usage:
  grpc-client <command> [flags]

Commands:
  get-me        Get current user by JWT
  get-user      Get user by ID
  checkout      Checkout active cart
  list-orders   List current user's orders
  get-order     Get order by ID
  watch-order   Watch order status stream

Examples:
  grpc-client get-me --token "ACCESS_TOKEN"
  grpc-client get-user --user-id "USER_ID"
  grpc-client checkout --token "ACCESS_TOKEN"
  grpc-client list-orders --token "ACCESS_TOKEN"
  grpc-client get-order --token "ACCESS_TOKEN" --order-id "ORDER_ID"
  grpc-client watch-order --token "ACCESS_TOKEN" --order-id "ORDER_ID"
`)
}
