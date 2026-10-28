package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDynamicSchemaRouter(t *testing.T) {
	router := DynamicSchemaRouter("schema_core", nil)

	handler := router(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activeSchema := GetActiveSchema(r.Context())
		w.Header().Set("X-Active-Schema", activeSchema)
		w.WriteHeader(http.StatusOK)
	}))

	// Case 1: Default fallback
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if s := rr1.Header().Get("X-Active-Schema"); s != "schema_core" {
		t.Errorf("expected schema_core, got %s", s)
	}

	// Case 2: URL prefix /api/v1/menupulse -> schema_menupulse
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/menupulse/items", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if s := rr2.Header().Get("X-Active-Schema"); s != "schema_menupulse" {
		t.Errorf("expected schema_menupulse, got %s", s)
	}

	// Case 3: Explicit header X-ToolGrid-Schema
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/custom", nil)
	req3.Header.Set("X-ToolGrid-Schema", "schema_custom_tenant")
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	if s := rr3.Header().Get("X-Active-Schema"); s != "schema_custom_tenant" {
		t.Errorf("expected schema_custom_tenant, got %s", s)
	}
}
