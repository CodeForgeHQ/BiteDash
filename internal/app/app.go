package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bitedash/internal/config"
	authpkg "bitedash/internal/pkg/auth"

	"google.golang.org/grpc"
)

type App struct {
	cfg          *config.Config
	diContainer  *diContainer
	httpServer   *http.Server
	grpcServer   *grpc.Server
	grpcListener net.Listener
}

func Build(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	configureLogger(cfg.LogLevel)
	authpkg.SetJWTSecret(cfg.JWTSecret)

	container, err := newDIContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize dependencies: %w", err)
	}

	a := &App{cfg: cfg, diContainer: container}
	a.httpServer = &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      container.httpServer.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	a.grpcServer = container.grpcServer.Server()
	return a, nil
}

func configureLogger(rawLevel string) {
	level := new(slog.LevelVar)
	level.Set(parseLogLevel(rawLevel))
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func parseLogLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (a *App) Run() error {
	listener, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
	if err != nil {
		_ = a.diContainer.Close()
		return fmt.Errorf("listen on grpc port %s: %w", a.cfg.GRPCPort, err)
	}
	a.grpcListener = listener

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	go func() {
		slog.Info("http server started", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()
	go func() {
		slog.Info("grpc server started", "addr", listener.Addr().String())
		if err := a.grpcServer.Serve(listener); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return a.shutdown()
	case runErr := <-errCh:
		return errors.Join(runErr, a.shutdown())
	}
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpErr := a.httpServer.Shutdown(shutdownCtx)
	grpcStopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(grpcStopped)
	}()
	select {
	case <-grpcStopped:
		slog.Info("grpc server stopped gracefully")
	case <-shutdownCtx.Done():
		slog.Warn("grpc graceful stop timeout, forcing stop")
		a.grpcServer.Stop()
	}

	dbErr := a.diContainer.Close()
	if err := errors.Join(httpErr, dbErr); err != nil {
		return fmt.Errorf("shutdown application: %w", err)
	}
	slog.Info("application stopped gracefully")
	return nil
}
