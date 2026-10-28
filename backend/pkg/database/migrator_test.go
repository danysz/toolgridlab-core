package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRunCoreMigrations_InMemoryFallback(t *testing.T) {
	ctx := context.Background()
	db, err := Connect(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error on empty DATABASE_URL: %v", err)
	}
	if err := db.RunCoreMigrations(ctx); err != nil {
		t.Fatalf("expected in-memory migration no-op to succeed, got: %v", err)
	}
}

func TestRunCoreMigrations_LiveDatabase(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set; skipping live migration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := Connect(ctx, dbURL)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer db.Close()

	// First run: applies unapplied migrations
	if err := db.RunCoreMigrations(ctx); err != nil {
		t.Fatalf("first RunCoreMigrations failed: %v", err)
	}

	// Second run: verifies idempotency and checksum matching
	if err := db.RunCoreMigrations(ctx); err != nil {
		t.Fatalf("second RunCoreMigrations (idempotency check) failed: %v", err)
	}

	records, err := db.AppliedMigrations(ctx)
	if err != nil {
		t.Fatalf("AppliedMigrations failed: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("expected at least 1 applied migration, got 0")
	}
	for _, r := range records {
		t.Logf("Verified migration: module=%s version=%s checksum=%s applied_at=%s",
			r.Module, r.Version, r.Checksum[:12], r.AppliedAt.Format(time.RFC3339))
	}
}
