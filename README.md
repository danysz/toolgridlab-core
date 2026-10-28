# ToolGridLab Core Engine (`toolgridlab-core`)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Flutter Version](https://img.shields.io/badge/Flutter-3.20+-02569B?style=flat&logo=flutter)](https://flutter.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Architecture: Modular Monolith](https://img.shields.io/badge/Architecture-Modular%20Monolith-purple)](https://toolgridlab.com)

The foundational open-source core engine and boilerplate for building scalable micro-SaaS portfolios. Designed for hyper-specific utilities with single centralized authentication, isolated database schemas per tool, and zero idle baseline cloud costs.

---

## 🏛️ Core Architecture Pillars

### 1. Dynamic PostgreSQL Schema Routing
* All micro-utilities share a single centralized Supabase PostgreSQL instance while maintaining strict data isolation.
* Incoming requests dynamically switch search paths (`SET search_path TO schema_<tool>;`) based on route prefixes or tenant headers.
* **Spin-Off Ready:** Any individual tool schema can be exported and pointed to an independent Go binary with zero vendor lock-in.

### 2. Go Modular Monolith Backend
* Built with Go 1.22+ on Google Cloud Run (containerized execution with zero idle baseline cost).
* Strict route versioning (`/api/v1/...`).
* Generic pluggable module interface (`core.Module`) allowing commercial micro-tools to register isolated endpoints and schemas cleanly.

### 3. Decoupled Supabase OIDC Auth
* Unified Google & Apple OAuth sign-in via OpenID Connect (OIDC).
* Centralized user profile pointer (`schema_core.user_profiles`) mapping external Supabase Auth UIDs to internal application data with zero local password storage.

### 4. Edge-Cached Flutter Web PWA
* High-performance Flutter Web application optimized for Cloudflare Pages edge deployment.
* White-label configuration loader (`AppConfig`) for domain-driven customizations.
* Unified Sentry error monitoring across both frontend and backend.

---

## 📂 Repository Structure

```text
toolgridlab-core/
├── backend/
│   ├── cmd/server/main.go            # Monolith entrypoint & graceful shutdown
│   ├── internal/
│   │   ├── config/config.go          # Environment configuration loader
│   │   ├── core/                     # Pluggable module registry & user profiles
│   │   │   ├── module.go             # Module interface definition
│   │   │   ├── handler.go            # Healthz, email verifier, and module discovery
│   │   │   └── user_service.go       # User profile synchronization
│   │   ├── database/                 # Supabase PostgreSQL pool & migrations
│   │   │   ├── postgres.go           # pgxpool connection manager
│   │   │   └── migrations/           # schema_core initialization
│   │   ├── middleware/               # HTTP middleware pipeline
│   │   │   ├── schema_router.go      # Dynamic search_path routing
│   │   │   ├── auth.go               # Supabase OIDC JWT verification
│   │   │   ├── cors.go               # Cloudflare Pages CORS handling
│   │   │   ├── logging.go            # Structured HTTP logging
│   │   │   └── sentry.go             # Sentry panic recovery
│   │   └── verification/             # Reusable verification utilities
│   │       ├── email_verifier.go     # RFC 5322 regex, disposable block, DNS MX lookup
│   │       └── ip_anonymizer.go      # SHA-256 IP anonymization
│   ├── Dockerfile                    # Multi-stage distroless build for Cloud Run
│   └── .env.example
├── frontend/
│   ├── lib/
│   │   ├── main.dart                 # App initialization (Sentry + Supabase)
│   │   ├── config/app_config.dart    # White-label configuration loader
│   │   ├── auth/auth_service.dart    # Google & Apple OIDC OAuth service
│   │   ├── api/core_api_client.dart  # Core HTTP client with schema headers
│   │   ├── theme/app_theme.dart      # Dark mode design tokens
│   │   ├── screens/home_screen.dart  # Interactive starter dashboard
│   │   └── widgets/                  # Navigation bar, auth modal, cards, footer
│   ├── build.sh                      # Cloudflare Pages release build script
│   └── wrangler.toml                 # Cloudflare Pages configuration
├── LICENSE                           # MIT License
└── .gitignore
```

---

## 🚀 Getting Started

### 1. Prerequisites
- **Go**: 1.22 or newer
- **Flutter**: 3.20 or newer
- **Docker** (optional, for container runs)

### 2. Backend Setup

```bash
cd backend

# Copy sample environment configuration
cp .env.example .env

# Run database migration (Supabase SQL Editor or psql)
# See: internal/database/migrations/001_core_schema.sql

# Run tests
go test -v ./...

# Start server
go run ./cmd/server
```

The server will listen on `http://localhost:8080` with the following core endpoints:
- `GET  /api/v1/healthz` — Health check & uptime
- `POST /api/v1/core/verify-email` — RFC 5322 regex + disposable domain + DNS MX check
- `GET  /api/v1/core/modules` — List registered pluggable modules
- `GET  /api/v1/core/me` — Authenticated user profile (requires Supabase JWT)

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
flutter pub get

# Run tests
flutter test

# Run Flutter Web locally
flutter run -d chrome
```

---

## 🧩 Adding a New Micro-Tool Module

To add an isolated micro-utility module from your private repository to this core engine:

1. Create a dedicated schema (e.g. `schema_mytool`).
2. Implement the `core.Module` interface:

```go
package mytool

import (
    "net/http"
    "github.com/toolgridlab/toolgridlab-core/backend/internal/core"
)

type MyToolModule struct{}

func (m *MyToolModule) ID() string     { return "mytool" }
func (m *MyToolModule) Schema() string { return "schema_mytool" }

func (m *MyToolModule) RegisterRoutes(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
    // Public endpoint (automatically routed to schema_mytool)
    mux.HandleFunc("GET /api/v1/mytool/items", handleGetItems)

    // Protected endpoint requiring Supabase OIDC Auth
    mux.Handle("POST /api/v1/mytool/items", requireAuth(http.HandlerFunc(handleCreateItem)))
}
```

3. Register your module in `cmd/server/main.go`:

```go
moduleRegistry.Register(&mytool.MyToolModule{})
```

The dynamic schema router will automatically intercept any request matching `/api/v1/mytool/*` and route queries to `schema_mytool`.

---

## 🚢 Production Deployment

### Backend (Google Cloud Run)
```bash
docker build -t gcr.io/[PROJECT_ID]/toolgridlab-core:latest ./backend
gcloud run deploy toolgridlab-core \
  --image gcr.io/[PROJECT_ID]/toolgridlab-core:latest \
  --platform managed \
  --allow-unauthenticated \
  --port 8080
```

### Frontend (Cloudflare Pages)
* Build command: `./build.sh` (or `flutter build web --release`)
* Output directory: `frontend/build/web`

---

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
