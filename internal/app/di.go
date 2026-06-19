package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"bitedash/internal/config"
	db "bitedash/internal/db/sqlc"
	"bitedash/internal/grpcserver"
	grpchandler "bitedash/internal/grpcserver/handler"
	"bitedash/internal/handler"
	"bitedash/internal/server"
	"bitedash/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// diContainer owns the application's long-lived dependencies. Construction is
// eager so startup errors are returned instead of terminating the process from
// deep inside the dependency graph.
type diContainer struct {
	db         *sql.DB
	httpServer *server.Server
	grpcServer *grpcserver.Server
}

func newDIContainer(ctx context.Context, cfg *config.Config) (*diContainer, error) {
	pool, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.PingContext(pingCtx); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	queries := db.New(pool)
	authService := service.NewAuthService(queries)
	restaurantService := service.NewRestaurantService(queries)
	cartService := service.NewCartService(queries, pool)
	orderService := service.NewOrderService(pool, queries)
	productService := service.NewProductService(queries)

	httpServer := server.NewServer(server.Deps{
		DB:                pool,
		AuthHandler:       handler.NewAuthHandler(authService),
		RestaurantHandler: handler.NewRestaurantHandler(restaurantService),
		CartHandler:       handler.NewCartHandler(cartService),
		OrderHandler:      handler.NewOrderHandler(orderService),
		ProductHandler:    handler.NewProductHandler(productService),
	})
	grpcServer := grpcserver.NewServer(grpcserver.Deps{
		UserHandler:  grpchandler.NewUserHandler(authService),
		OrderHandler: grpchandler.NewOrderHandler(orderService),
	})

	return &diContainer{db: pool, httpServer: httpServer, grpcServer: grpcServer}, nil
}

func (c *diContainer) Close() error {
	if c == nil || c.db == nil {
		return nil
	}
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}
