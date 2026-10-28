-- ToolGridLab: Core Schema & Global Pointer Architecture
-- Schema: schema_core

CREATE SCHEMA IF NOT EXISTS schema_core;

SET search_path TO schema_core;

-- 1. Decoupled Central User Profiles (mapped to external Supabase Auth IDs)
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id TEXT NOT NULL UNIQUE,       -- Supabase Auth UID (from Google / Apple OIDC)
    email TEXT NOT NULL,
    provider TEXT NOT NULL DEFAULT 'oidc',  -- 'google', 'apple', etc.
    display_name TEXT,
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_email ON user_profiles(email);

-- 2. Tenant & Micro-Tool Schema Registry (Dynamic Schema Pointer)
CREATE TABLE IF NOT EXISTS schema_registry (
    id TEXT PRIMARY KEY,                    -- e.g. 'validation_engine', 'menupulse'
    schema_name TEXT NOT NULL UNIQUE,       -- e.g. 'schema_validation_engine', 'schema_menupulse'
    display_name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Register default schemas
INSERT INTO schema_registry (id, schema_name, display_name, description)
VALUES 
    ('core', 'schema_core', 'Core Engine & Global Auth', 'Central authentication and user profile registry')
ON CONFLICT (id) DO NOTHING;
