package engine

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/config"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/core"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/database"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/middleware"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/verification"
)

// Run boots the ToolGridLab Core Engine and registers any supplied pluggable modules.
func Run(modules ...core.Module) {
	cfg := config.Load()
	log.Printf("Starting ToolGridLab Core Engine (env=%s, port=%s)...", cfg.Environment, cfg.Port)

	// 1. Initialize Sentry (if configured)
	if cfg.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			Environment:      cfg.Environment,
			TracesSampleRate: 0.2,
		})
		if err != nil {
			log.Printf("[SENTRY] Initialization failed: %v", err)
		} else {
			defer sentry.Flush(2 * time.Second)
			log.Println("[SENTRY] Sentry initialized successfully.")
		}
	}

	// 2. Connect to Supabase PostgreSQL & Run Embedded Core Migrations
	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[DB] Fatal connection error: %v", err)
	}
	defer db.Close()

	if err := db.RunCoreMigrations(ctx); err != nil {
		log.Fatalf("[DB] Core schema migration failed: %v", err)
	}

	// 3. Initialize Core Services & Module Registry
	userService := core.NewUserService(db.Pool)
	emailVerifier := verification.NewEmailVerifier(3 * time.Second)
	moduleRegistry := core.NewModuleRegistry()

	for _, m := range modules {
		if err := moduleRegistry.Register(m); err != nil {
			log.Fatalf("[MODULE] Registration failed: %v", err)
		}
	}

	// 4. Setup ServeMux & Core Handlers
	mux := http.NewServeMux()
	requireAuth := middleware.SupabaseAuthMiddleware(cfg.SupabaseJWTSecret, false)

	coreHandler := core.NewCoreHandler(userService, emailVerifier, moduleRegistry)
	coreHandler.RegisterRoutes(mux, requireAuth)

	// Initialize, migrate, and mount routes for all registered modules
	for _, m := range moduleRegistry.All() {
		if configurable, ok := m.(core.ConfigurableModule); ok {
			if err := configurable.Init(db, cfg); err != nil {
				log.Fatalf("[MODULE] Init failed for module '%s': %v", m.ID(), err)
			}
		}
		if migratable, ok := m.(core.MigratableModule); ok {
			fsys, dir := migratable.Migrations()
			if err := db.RunMigrationsFS(ctx, m.ID(), fsys, dir); err != nil {
				log.Fatalf("[DB] Migration failed for module '%s': %v", m.ID(), err)
			}
		}
		m.RegisterRoutes(mux, requireAuth)
	}

	// 5. Wrap with Global Middleware Pipeline
	var handler http.Handler = mux
	handler = middleware.DynamicSchemaRouter(middleware.DefaultSchema, nil)(handler)
	handler = middleware.SentryRecovery(handler)
	handler = middleware.Logging(handler)
	handler = middleware.CORS(cfg.AllowedCORS)(handler)

	// 6. Graceful HTTP Server Execution
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("ToolGridLab Core Engine listening on http://0.0.0.0:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped cleanly.")
}

