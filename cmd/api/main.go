package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	db "bitedash/internal/db/query"
	"bitedash/internal/server"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/food_delivery?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("db ping failed:", err)
	}

	log.Println("Connected to database")

	queries := db.New(pool)

	srv := server.NewServer(pool, queries)

	log.Println("Starting server on :8080")

	if err := srv.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
