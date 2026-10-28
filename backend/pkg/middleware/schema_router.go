package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type contextKey string

const (
	SchemaContextKey contextKey = "db_search_path"
	DefaultSchema    string     = "schema_core"
)

var validSchemaName = regexp.MustCompile(`^[a-z0-9_]{1,63}$`)

// SchemaResolverFunc allows modules or apps to provide custom schema resolution logic
type SchemaResolverFunc func(r *http.Request) string

// DynamicSchemaRouter dynamically sets the active PostgreSQL schema in request context.
// It resolves the target schema via:
// 1. Explicit header 'X-ToolGrid-Schema'
// 2. Custom SchemaResolverFunc (if provided)
// 3. First URL segment under /api/v1/<module>/ -> schema_<module>
// 4. Default schema fallback (e.g., schema_core)
func DynamicSchemaRouter(defaultSchema string, customResolver SchemaResolverFunc) func(http.Handler) http.Handler {
	if defaultSchema == "" {
		defaultSchema = DefaultSchema
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			targetSchema := ""

			// 1. Check explicit header
			if headerSchema := r.Header.Get("X-ToolGrid-Schema"); headerSchema != "" {
				targetSchema = sanitizeSchemaName(headerSchema)
			}

			// 2. Check custom resolver
			if targetSchema == "" && customResolver != nil {
				if resolved := customResolver(r); resolved != "" {
					targetSchema = sanitizeSchemaName(resolved)
				}
			}

			// 3. URL path-based schema switching (/api/v1/<module>/...)
			if targetSchema == "" {
				path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
				parts := strings.Split(path, "/")
				if len(parts) > 0 && parts[0] != "" && parts[0] != "healthz" && parts[0] != "auth" {
					targetSchema = "schema_" + sanitizeSchemaName(parts[0])
				}
			}

			// 4. Fallback default
			if targetSchema == "" {
				targetSchema = defaultSchema
			}

			ctx := context.WithValue(r.Context(), SchemaContextKey, targetSchema)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetActiveSchema extracts the schema from context
func GetActiveSchema(ctx context.Context) string {
	if schema, ok := ctx.Value(SchemaContextKey).(string); ok && schema != "" {
		return schema
	}
	return DefaultSchema
}

func sanitizeSchemaName(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	var clean strings.Builder
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			clean.WriteRune(r)
		}
	}
	res := clean.String()
	if !validSchemaName.MatchString(res) {
		return DefaultSchema
	}
	return res
}
