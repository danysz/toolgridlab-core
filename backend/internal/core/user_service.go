package core

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProfile struct {
	ID          string    `json:"id"`
	AuthUserID  string    `json:"auth_user_id"`
	Email       string    `json:"email"`
	Provider    string    `json:"provider"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{pool: pool}
}

// SyncUserProfile upserts the internal user profile upon OIDC sign-in
func (s *UserService) SyncUserProfile(ctx context.Context, authUserID, email, provider, displayName string) (*UserProfile, error) {
	if s.pool == nil {
		// Mock profile for local/testing without database
		return &UserProfile{
			ID:          "mock-user-id",
			AuthUserID:  authUserID,
			Email:       email,
			Provider:    provider,
			DisplayName: displayName,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	query := `
		INSERT INTO schema_core.user_profiles (auth_user_id, email, provider, display_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (auth_user_id) DO UPDATE SET
			email = EXCLUDED.email,
			provider = EXCLUDED.provider,
			display_name = COALESCE(EXCLUDED.display_name, schema_core.user_profiles.display_name),
			updated_at = NOW()
		RETURNING id, auth_user_id, email, provider, display_name, COALESCE(avatar_url, ''), is_active, created_at, updated_at;
	`

	var u UserProfile
	err := s.pool.QueryRow(ctx, query, authUserID, email, provider, displayName).Scan(
		&u.ID,
		&u.AuthUserID,
		&u.Email,
		&u.Provider,
		&u.DisplayName,
		&u.AvatarURL,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user profile: %w", err)
	}

	return &u, nil
}
