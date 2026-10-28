#!/bin/bash
set -e

# Default to production if not specified by Cloudflare environment
APP_ENV=${APP_ENV:-"production"}
APP_NAME=${APP_NAME:-"ToolGridLab"}

echo "=== Building ToolGridLab Frontend for Cloudflare Pages (Environment: $APP_ENV) ==="

# Default API URL based on environment
if [ "$APP_ENV" = "staging" ] || [ "$APP_ENV" = "test" ]; then
  API_BASE_URL=${API_BASE_URL:-"https://api-staging.toolgridlab.com"}
else
  API_BASE_URL=${API_BASE_URL:-"https://api.toolgridlab.com"}
fi

BUILD_ARGS=(
  "web"
  "--release"
  "--dart-define=APP_ENV=$APP_ENV"
  "--dart-define=APP_NAME=$APP_NAME"
  "--dart-define=API_BASE_URL=$API_BASE_URL"
)

# Optional Supabase and Sentry injections from Cloudflare Pages Environment Variables
if [ -n "$SUPABASE_URL" ]; then
  BUILD_ARGS+=("--dart-define=SUPABASE_URL=$SUPABASE_URL")
fi

if [ -n "$SUPABASE_ANON_KEY" ]; then
  BUILD_ARGS+=("--dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY")
fi

if [ -n "$SENTRY_DSN" ]; then
  BUILD_ARGS+=("--dart-define=SENTRY_DSN=$SENTRY_DSN")
fi

flutter build "${BUILD_ARGS[@]}"

echo "=== Build finished successfully! Output in build/web ==="
