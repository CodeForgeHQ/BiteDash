package main

import (
	_ "bitedash/docs"
	"context"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Регистрация драйвера pgx для database/sql

	"bitedash/internal/app"
)

// @title BiteDash API
// @version 1.0
// @description Backend API for BiteDash food delivery service
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	a, err := app.Build(context.Background())
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}
