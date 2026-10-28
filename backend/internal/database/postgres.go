package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		log.Println("[DB] DATABASE_URL not set. Running with In-Memory fallback store.")
		return &DB{Pool: nil}, nil
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	// Ping database
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping Supabase PostgreSQL: %w", err)
	}

	log.Println("[DB] Connected to Supabase PostgreSQL successfully.")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// SetSchema executes SET search_path TO <schema>; on a connection
func (db *DB) SetSchema(ctx context.Context, schema string) error {
	if db.Pool == nil {
		return nil
	}
	_, err := db.Pool.Exec(ctx, fmt.Sprintf("SET search_path TO %s;", schema))
	return err
}
