#!/bin/bash
set -e

# Default to production if not specified by Cloudflare environment
# On Cloudflare Pages, CF_PAGES_BRANCH is set automatically (e.g., "staging" vs "main")
if [ -z "$APP_ENV" ]; then
  if [ "$CF_PAGES_BRANCH" = "staging" ]; then
    APP_ENV="staging"
  else
    APP_ENV="production"
  fi
fi

# Load local .env.<APP_ENV> file if present and variables aren't already set in environment
if [ -f ".env.$APP_ENV" ]; then
  set -a
  # shellcheck disable=SC1090
  source ".env.$APP_ENV"
  set +a
fi

APP_NAME=${APP_NAME:-"ToolGridLab"}

echo "=== Building ToolGridLab Frontend for Cloudflare Pages (Environment: $APP_ENV) ==="

# Bootstrap Flutter SDK on Cloudflare Pages CI if not already in PATH
if ! command -v flutter &> /dev/null; then
  echo "=== Flutter not found in PATH; bootstrapping Flutter stable SDK for CI ==="
  if [ ! -d "$HOME/flutter" ]; then
    git clone https://github.com/flutter/flutter.git -b stable --depth 1 "$HOME/flutter"
  fi
  export PATH="$HOME/flutter/bin:$PATH"
fi

flutter --version
flutter pub get

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
