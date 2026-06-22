package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/yallaops/yallaops/core/internal/api"
	"github.com/yallaops/yallaops/core/internal/release"
	"github.com/yallaops/yallaops/core/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Error("failed to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	log.Info("database connected")

	// stdlib adapter wraps pgxpool as database/sql.DB for sqlc compatibility
	db := stdlib.OpenDBFromPool(pool)
	defer func() { _ = db.Close() }()

	q := store.New(db)
	releaseSvc := release.NewService(q)
	releaseHandler := api.NewReleaseHandler(releaseSvc)

	srv := api.NewServer(releaseHandler)

	addr := ":50051"
	log.Info("starting gRPC server", "addr", addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe(addr)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gracefully")
		srv.GracefulStop()
	case err := <-errCh:
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}
