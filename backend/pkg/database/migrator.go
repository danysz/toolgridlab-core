package database

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var coreMigrationsFS embed.FS

// advisoryLockID is a fixed key used with pg_advisory_lock to serialize migrations across Cloud Run instances.
const advisoryLockID int64 = 84665271901

// MigrationRecord represents an applied database migration in schema_core.schema_migrations.
type MigrationRecord struct {
	Module    string    `json:"module"`
	Version   string    `json:"version"`
	Checksum  string    `json:"checksum"`
	AppliedAt time.Time `json:"applied_at"`
}

// RunCoreMigrations executes all embedded core SQL migrations in migrations/*.sql.
func (db *DB) RunCoreMigrations(ctx context.Context) error {
	return db.RunMigrationsFS(ctx, "core", coreMigrationsFS, "migrations")
}

// RunMigrationsFS executes all *.sql files in dir from fsys for a given module in lexicographical order,
// tracking applied versions and SHA-256 checksums in schema_core.schema_migrations.
func (db *DB) RunMigrationsFS(ctx context.Context, module string, fsys fs.FS, dir string) error {
	if db.Pool == nil {
		log.Printf("[MIGRATOR] Skipping migrations for module '%s' (in-memory mode).", module)
		return nil
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory '%s' for module '%s': %w", dir, module, err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Printf("[MIGRATOR] No SQL migration files found for module '%s'.", module)
		return nil
	}

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire DB connection for migrations: %w", err)
	}
	defer conn.Release()

	// Serialize migrations across concurrent instances
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", advisoryLockID); err != nil {
		return fmt.Errorf("failed to acquire migration advisory lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = conn.Exec(unlockCtx, "SELECT pg_advisory_unlock($1)", advisoryLockID)
	}()

	// Ensure schema_core and schema_migrations version table exist
	initSQL := `
		CREATE SCHEMA IF NOT EXISTS schema_core;
		CREATE TABLE IF NOT EXISTS schema_core.schema_migrations (
			module TEXT NOT NULL,
			version TEXT NOT NULL,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (module, version)
		);
	`
	if _, err := conn.Exec(ctx, initSQL); err != nil {
		return fmt.Errorf("failed to initialize schema_core.schema_migrations: %w", err)
	}

	appliedCount := 0
	for _, filename := range sqlFiles {
		filePath := path.Join(dir, filename)
		contentBytes, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file '%s': %w", filePath, err)
		}

		hash := sha256.Sum256(contentBytes)
		checksum := hex.EncodeToString(hash[:])

		var existingChecksum string
		err = conn.QueryRow(
			ctx,
			"SELECT checksum FROM schema_core.schema_migrations WHERE module = $1 AND version = $2",
			module, filename,
		).Scan(&existingChecksum)

		if err == nil {
			if existingChecksum != checksum {
				return fmt.Errorf(
					"checksum mismatch for already-applied migration %s/%s (expected %s, got %s)",
					module, filename, existingChecksum, checksum,
				)
			}
			continue
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to query migration status for %s/%s: %w", module, filename, err)
		}

		// Execute migration and version tracking insert atomically in a transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s/%s: %w", module, filename, err)
		}

		if _, err := tx.Exec(ctx, string(contentBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s/%s: %w", module, filename, err)
		}

		if _, err := tx.Exec(
			ctx,
			"INSERT INTO schema_core.schema_migrations (module, version, checksum) VALUES ($1, $2, $3)",
			module, filename, checksum,
		); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %s/%s: %w", module, filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s/%s: %w", module, filename, err)
		}

		appliedCount++
		log.Printf("[MIGRATOR] Applied migration %s/%s (checksum=%s)", module, filename, checksum[:12])
	}

	if appliedCount == 0 {
		log.Printf("[MIGRATOR] Module '%s' database schema is up to date (%d migration(s) verified).", module, len(sqlFiles))
	} else {
		log.Printf("[MIGRATOR] Module '%s' applied %d new migration(s) (%d total).", module, appliedCount, len(sqlFiles))
	}

	return nil
}

// AppliedMigrations returns all recorded migrations ordered by applied_at.
func (db *DB) AppliedMigrations(ctx context.Context) ([]MigrationRecord, error) {
	if db.Pool == nil {
		return nil, nil
	}

	rows, err := db.Pool.Query(
		ctx,
		"SELECT module, version, checksum, applied_at FROM schema_core.schema_migrations ORDER BY applied_at ASC, module ASC, version ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list applied migrations: %w", err)
	}
	defer rows.Close()

	var records []MigrationRecord
	for rows.Next() {
		var r MigrationRecord
		if err := rows.Scan(&r.Module, &r.Version, &r.Checksum, &r.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration record: %w", err)
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

