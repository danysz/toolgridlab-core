package core

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toolgridlab/toolgridlab-core/backend/pkg/middleware"
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/verification"
)

type CoreHandler struct {
	userService   *UserService
	emailVerifier *verification.EmailVerifier
	registry      *ModuleRegistry
	startTime     time.Time
}

func NewCoreHandler(us *UserService, ev *verification.EmailVerifier, reg *ModuleRegistry) *CoreHandler {
	return &CoreHandler{
		userService:   us,
		emailVerifier: ev,
		registry:      reg,
		startTime:     time.Now(),
	}
}

func (h *CoreHandler) RegisterRoutes(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	// Public Healthz
	mux.HandleFunc("GET /api/v1/healthz", h.HandleHealthz)

	// Generic Reusable Email Verification Utility
	mux.HandleFunc("POST /api/v1/core/verify-email", h.HandleVerifyEmail)

	// Zero-latency Edge Geo / Country Detection
	mux.HandleFunc("GET /api/v1/core/geo", h.HandleGeo)

	// List Registered Modules
	mux.HandleFunc("GET /api/v1/core/modules", h.HandleListModules)

	// Authenticated User Profile
	mux.Handle("GET /api/v1/core/me", requireAuth(http.HandlerFunc(h.HandleGetMe)))
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Service   string    `json:"service"`
}

func (h *CoreHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	middleware.RespondJSON(w, http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Service:   "toolgridlab-core",
	})
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
}

type VerifyEmailResponse struct {
	Valid   bool   `json:"valid"`
	Email   string `json:"email"`
	Country string `json:"country,omitempty"`
	Message string `json:"message"`
}

func (h *CoreHandler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	country := verification.ExtractClientCountry(r, "US")

	if err := h.emailVerifier.Verify(r.Context(), req.Email); err != nil {
		middleware.RespondJSON(w, http.StatusOK, VerifyEmailResponse{
			Valid:   false,
			Email:   req.Email,
			Country: country,
			Message: err.Error(),
		})
		return
	}

	middleware.RespondJSON(w, http.StatusOK, VerifyEmailResponse{
		Valid:   true,
		Email:   req.Email,
		Country: country,
		Message: "email is valid and reachable",
	})
}

func (h *CoreHandler) HandleGeo(w http.ResponseWriter, r *http.Request) {
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"country": verification.ExtractClientCountry(r, "US"),
	})
}

func (h *CoreHandler) HandleListModules(w http.ResponseWriter, r *http.Request) {
	type ModuleInfo struct {
		ID     string `json:"id"`
		Schema string `json:"schema"`
	}

	modules := h.registry.All()
	res := make([]ModuleInfo, len(modules))
	for i, m := range modules {
		res[i] = ModuleInfo{
			ID:     m.ID(),
			Schema: m.Schema(),
		}
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"modules": res,
		"count":   len(res),
	})
}

func (h *CoreHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		middleware.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	profile, err := h.userService.SyncUserProfile(r.Context(), authUser.ID, authUser.Email, authUser.Provider, "")
	if err != nil {
		middleware.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	middleware.RespondJSON(w, http.StatusOK, profile)
}
