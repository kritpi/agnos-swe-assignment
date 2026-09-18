package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kritpi/agnos-swe-assignment/handler"
	"github.com/kritpi/agnos-swe-assignment/internal/adapter"
	"github.com/kritpi/agnos-swe-assignment/internal/core/service"
	"github.com/kritpi/agnos-swe-assignment/property"
	"github.com/kritpi/agnos-swe-assignment/repository"
	"github.com/kritpi/agnos-swe-assignment/router"
)

// Run loads config, wires dependencies, and starts the HTTP server.
func Run() error {
	ctx := context.Background()

	cfg, err := property.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// PostgreSQL connection pool.
	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("init postgres pool: %w", err)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	repo := repository.New(pool, cfg.DBTable)
	adapters := adapter.New(cfg)
	svc := service.New(repo, adapters, cfg)
	h := handler.New(svc)

	r := gin.Default()
	router.SetUpRouter(r, h)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}
