#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "[INSTALL] Missing dependency: $1"
    exit 1
  }
}

ask_required() {
  local key="$1"
  local prompt="$2"
  local value=""
  while [[ -z "${value}" ]]; do
    read -r -p "${prompt}: " value
    value="${value//[$'\r\n']}"
  done
  printf '%s' "$value"
}

ask_optional() {
  local prompt="$1"
  local default="$2"
  local value=""
  read -r -p "${prompt} (default: ${default}): " value
  value="${value//[$'\r\n']}"
  printf '%s' "${value:-$default}"
}

ask_optional_secret() {
  local prompt="$1"
  local default="$2"
  local value=""
  read -r -s -p "${prompt} (default: ${default}): " value; echo
  value="${value//[$'\r\n']}"
  printf '%s' "${value:-$default}"
}

echo "[INSTALL] Requesting sudo rights..."
sudo -v

echo "[INSTALL] Checking dependencies..."
require_cmd podman

ENV_FILE=".env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "[INSTALL] No .env found at ${PROJECT_DIR}/.env"
  read -r -p "[INSTALL] Create it interactively now? (y/n): " yn
  if [[ "${yn,,}" != "y" ]]; then
    echo "[INSTALL] Aborting (user refused interactive setup)."
    exit 1
  fi

  FRONTEND_PORT="$(ask_optional "FRONTEND_PORT" "20301")"
  BACKEND_PORT="$(ask_optional "BACKEND_PORT" "9001")"
  PUBLIC_API_URL="$(ask_required "PUBLIC_API_URL" "PUBLIC_API_URL (required, e.g. https://lm.example.com or http://localhost:PORT)")"

  JWT_SECRET="$(ask_required "JWT_SECRET" "JWT_SECRET (required)")"

  DATA_DIR="$(ask_optional "DATA_DIR" "./backend/data")"
  DB_PATH="$(ask_optional "DB_PATH" "/data/nerdy_link_manager.db")"

  cat > "$ENV_FILE" <<EOF
# Frontend build-time configuration
PUBLIC_API_URL=${PUBLIC_API_URL}

# Backend configuration
JWT_SECRET=${JWT_SECRET}
DB_PATH=${DB_PATH}

# Infra / host ports
FRONTEND_PORT=${FRONTEND_PORT}
BACKEND_PORT=${BACKEND_PORT}

# Persistence (mounted into pod as /data)
DATA_DIR=${DATA_DIR}
EOF

  chmod 600 "$ENV_FILE"
  echo "[INSTALL] .env created."
fi

# Load .env into this shell
set -a
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a

# Defaults if user created .env manually but forgot some vars
FRONTEND_PORT="${FRONTEND_PORT:-20301}"
BACKEND_PORT="${BACKEND_PORT:-9001}"
DATA_DIR="${DATA_DIR:-./backend/data}"
DB_PATH="${DB_PATH:-/data/nerdy_link_manager.db}"

mkdir -p "$DATA_DIR"

# Decide whether this is a first run (DB missing)
# If DB_PATH is under /data, map it to the host DATA_DIR.
DB_BASENAME="$(basename "$DB_PATH")"
HOST_DB_PATH="${DATA_DIR%/}/${DB_BASENAME}"

FIRST_RUN="false"
if [[ ! -f "$HOST_DB_PATH" ]]; then
  FIRST_RUN="true"
fi

if [[ "$FIRST_RUN" == "true" ]]; then
  echo "[INSTALL] No database found at ${HOST_DB_PATH} (first run detected)."

  # If not already set, ask for initial admin credentials and persist them in .env
  if [[ -z "${INITIAL_ADMIN_EMAIL:-}" ]]; then
    INITIAL_ADMIN_EMAIL="$(ask_required "INITIAL_ADMIN_EMAIL" "INITIAL_ADMIN_EMAIL (required for first run)")"
  fi

  if [[ -z "${INITIAL_ADMIN_PASSWORD:-}" ]]; then
    read -r -s -p "INITIAL_ADMIN_PASSWORD (required for first run): " INITIAL_ADMIN_PASSWORD; echo
    [[ -n "${INITIAL_ADMIN_PASSWORD}" ]] || { echo "[INSTALL] INITIAL_ADMIN_PASSWORD is required."; exit 1; }
  fi

  # Append to .env if missing (so next runs are consistent)
  if ! grep -q '^INITIAL_ADMIN_EMAIL=' "$ENV_FILE"; then
    echo "INITIAL_ADMIN_EMAIL=${INITIAL_ADMIN_EMAIL}" >> "$ENV_FILE"
  fi
  if ! grep -q '^INITIAL_ADMIN_PASSWORD=' "$ENV_FILE"; then
    echo "INITIAL_ADMIN_PASSWORD=${INITIAL_ADMIN_PASSWORD}" >> "$ENV_FILE"
  fi
  chmod 600 "$ENV_FILE"
else
  echo "[INSTALL] Existing database found at ${HOST_DB_PATH}."
fi

echo "[INSTALL] Building images..."
podman build -t nerdy-lm-backend -f backend/Dockerfile backend
podman build -t nerdy-lm-frontend --build-arg "PUBLIC_API_URL=${PUBLIC_API_URL}" -f frontend/Dockerfile frontend

echo "[INSTALL] Creating pod (if missing)..."
podman pod exists nerdy-lm-pod || podman pod create \
  --name nerdy-lm-pod \
  -p "0.0.0.0:${BACKEND_PORT}:8080" \
  -p "0.0.0.0:${FRONTEND_PORT}:80" \
  -v "${DATA_DIR}:/data:Z"

BACKEND_ENV_FILE=".nerdy-lm-backend.env"
cat > "$BACKEND_ENV_FILE" <<EOF
PORT=8080
JWT_SECRET=${JWT_SECRET}
DB_PATH=${DB_PATH}
EOF

# Only include initial admin envs when present (first run or user set them)
if [[ -n "${INITIAL_ADMIN_EMAIL:-}" ]]; then
  echo "INITIAL_ADMIN_EMAIL=${INITIAL_ADMIN_EMAIL}" >> "$BACKEND_ENV_FILE"
fi
if [[ -n "${INITIAL_ADMIN_PASSWORD:-}" ]]; then
  echo "INITIAL_ADMIN_PASSWORD=${INITIAL_ADMIN_PASSWORD}" >> "$BACKEND_ENV_FILE"
fi

chmod 600 "$BACKEND_ENV_FILE"

echo "[INSTALL] Starting containers..."
podman run --replace -d --pod nerdy-lm-pod --name nerdy-lm-backend \
  --env-file "$BACKEND_ENV_FILE" \
  nerdy-lm-backend:latest

podman run --replace -d --pod nerdy-lm-pod --name nerdy-lm-frontend \
  nerdy-lm-frontend:latest

echo "[INSTALL] Done."
echo "[INSTALL] Frontend: http://localhost:${FRONTEND_PORT}"
echo "[INSTALL] Backend:  http://localhost:${BACKEND_PORT}"
