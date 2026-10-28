package core

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/toolgridlab/toolgridlab-core/backend/pkg/config"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/database"
)

// Module represents a pluggable micro-utility in the modular monolith
type Module interface {
	// ID returns the module identifier (e.g., "validation", "menupulse")
	ID() string
	// Schema returns the dedicated PostgreSQL schema (e.g., "schema_validation")
	Schema() string
	// RegisterRoutes mounts the module's HTTP endpoints under /api/v1/<ID>/
	RegisterRoutes(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler)
}

// ConfigurableModule is an optional interface that modules can implement
// to receive the shared database handle and server configuration on startup.
type ConfigurableModule interface {
	Module
	Init(db *database.DB, cfg *config.Config) error
}

// MigratableModule is an optional interface that modules can implement
// to bundle embedded SQL migrations that are automatically executed and version-tracked on startup.
type MigratableModule interface {
	Module
	// Migrations returns the embedded filesystem and directory path containing *.sql files.
	Migrations() (fsys fs.FS, dir string)
}

// ModuleRegistry manages registered pluggable modules
type ModuleRegistry struct {
	modules map[string]Module
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
	}
}

func (r *ModuleRegistry) Register(m Module) error {
	id := m.ID()
	if _, exists := r.modules[id]; exists {
		return fmt.Errorf("module with ID '%s' is already registered", id)
	}
	r.modules[id] = m
	log.Printf("[MODULE] Registered module '%s' -> schema '%s'", id, m.Schema())
	return nil
}

func (r *ModuleRegistry) Get(id string) (Module, bool) {
	m, ok := r.modules[id]
	return m, ok
}

func (r *ModuleRegistry) All() []Module {
	list := make([]Module, 0, len(r.modules))
	for _, m := range r.modules {
		list = append(list, m)
	}
	return list
}
