package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"microservice-redis/internal/config"
	"microservice-redis/internal/storage/postgres"
	"time"
)

type Storage struct {
	*postgres.Queries
}

func Config() *pgxpool.Config {
	cfg := config.MustLoad().DB
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = cfg.MaxConns
	dbConfig.MinConns = cfg.MinConns
	dbConfig.MaxConnLifetime = time.Duration(cfg.MaxConnLifetime)
	dbConfig.MaxConnLifetime = time.Duration(cfg.MaxConnIdleTime)
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	return dbConfig
}

func New(log *slog.Logger) *Storage {
	conn, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Error("Failed to create a connection", "error", err)
		panic(err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Error("Failed to ping the database", "error", err)
	}

	queries := postgres.New(conn)

	return &Storage{queries}
}
