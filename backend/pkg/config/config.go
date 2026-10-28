package config

import (
	"os"
)

type Config struct {
	Port              string
	Environment       string
	DatabaseURL       string
	SentryDSN         string
	SupabaseURL       string
	SupabaseKey       string
	SupabaseJWTSecret string
	IPHashSalt        string
	AllowedCORS       string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	allowedCORS := os.Getenv("ALLOWED_CORS")
	if allowedCORS == "" {
		allowedCORS = "*"
	}

	salt := os.Getenv("IP_HASH_SALT")
	if salt == "" {
		salt = "toolgridlab-validation-salt-v1"
	}

	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SUPABASE_KEY")
	}

	return &Config{
		Port:              port,
		Environment:       env,
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		SentryDSN:         os.Getenv("SENTRY_DSN"),
		SupabaseURL:       os.Getenv("SUPABASE_URL"),
		SupabaseKey:       os.Getenv("SUPABASE_KEY"),
		SupabaseJWTSecret: jwtSecret,
		IPHashSalt:        salt,
		AllowedCORS:       allowedCORS,
	}
}
