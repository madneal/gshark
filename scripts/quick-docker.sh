#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

WITH_SCAN=false
ADMIN_USER="gshark"
ADMIN_PASSWORD="gshark"
# MySQL settings matching config.docker.yaml / compose
MYSQL_HOST="177.7.0.13"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD="madneal"
MYSQL_DB="gshark"
SKIP_INIT=false

usage() {
    cat <<'EOF'
Usage: scripts/quick-docker.sh [options]

Build and start GShark with Docker Compose, then initialize the database
if needed (admin account via flags — no browser required).

Options:
  --with-scan              Also start the scanner container.
  --admin-user NAME        Admin login username (default: gshark).
  --admin-password PASS    Admin login password (default: gshark).
  --skip-init              Do not run gshark init after start.
  -h, --help               Show this help message.

Example:
  ./scripts/quick-docker.sh --admin-user myadmin --admin-password 'S3cret!'
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --with-scan)
            WITH_SCAN=true
            ;;
        --admin-user)
            [[ $# -ge 2 ]] || { echo "--admin-user requires a value" >&2; exit 1; }
            ADMIN_USER="$2"
            shift
            ;;
        --admin-password)
            [[ $# -ge 2 ]] || { echo "--admin-password requires a value" >&2; exit 1; }
            ADMIN_PASSWORD="$2"
            shift
            ;;
        --skip-init)
            SKIP_INIT=true
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
    shift
done

command -v docker >/dev/null 2>&1 || {
    echo "Docker is required." >&2
    exit 1
}

docker info >/dev/null 2>&1 || {
    echo "Docker daemon is not reachable. Start Docker Desktop and retry." >&2
    exit 1
}

if docker compose version >/dev/null 2>&1; then
    COMPOSE=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose)
else
    echo "Docker Compose is required." >&2
    exit 1
fi

echo "[INFO] Building and starting mysql/server/web..."
"${COMPOSE[@]}" up -d --build mysql server web

if [[ "$WITH_SCAN" == true ]]; then
    echo "[INFO] Starting scan..."
    "${COMPOSE[@]}" up -d --build scan
fi

if [[ "$SKIP_INIT" != true ]]; then
    echo "[INFO] Waiting for MySQL and server..."
    for i in $(seq 1 60); do
        h=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' gshark-mysql 2>/dev/null || echo missing)
        if [[ "$h" == "healthy" ]] || [[ "$h" == "none" ]]; then
            if docker exec gshark-server ./gshark --help >/dev/null 2>&1 || \
               docker exec gshark-server ls ./gshark >/dev/null 2>&1; then
                break
            fi
        fi
        sleep 2
    done

    echo "[INFO] Initializing database (admin-user=${ADMIN_USER})..."
    # Run init inside the server container so it can reach MySQL on the compose network.
    # If already initialized, gshark init exits 0 with a skip message.
    if ! docker exec gshark-server ./gshark init \
        --host "$MYSQL_HOST" \
        --port "$MYSQL_PORT" \
        --user "$MYSQL_USER" \
        --password "$MYSQL_PASSWORD" \
        --db "$MYSQL_DB" \
        --admin-user "$ADMIN_USER" \
        --admin-password "$ADMIN_PASSWORD"; then
        echo "[WARN] gshark init failed; you can open http://localhost:8080 and initialize in the UI," >&2
        echo "       or re-run: docker exec gshark-server ./gshark init --help" >&2
    fi
fi

echo
"${COMPOSE[@]}" ps
echo
echo "GShark is starting at: http://localhost:8080"
echo "Admin login: ${ADMIN_USER} / (the password you set with --admin-password)"
