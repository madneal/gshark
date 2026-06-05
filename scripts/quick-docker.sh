#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

WITH_SCAN=false
BUILD=true
PULL=false

usage() {
    cat <<'EOF'
Usage: scripts/quick-docker.sh [options]

Quickly deploy GShark with Docker Compose.

Options:
  --with-scan     Start the scan container after mysql/server/web are up.
  --no-build      Start existing images without rebuilding.
  --pull          Pull newer base images before building.
  -h, --help      Show this help message.

Examples:
  scripts/quick-docker.sh
  scripts/quick-docker.sh --with-scan
EOF
}

log() {
    printf '[INFO] %s\n' "$*"
}

warn() {
    printf '[WARN] %s\n' "$*" >&2
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --with-scan)
            WITH_SCAN=true
            ;;
        --no-build)
            BUILD=false
            ;;
        --pull)
            PULL=true
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            fail "Unknown option: $1"
            ;;
    esac
    shift
done

command -v docker >/dev/null 2>&1 || fail "Docker is required."
docker info >/dev/null 2>&1 || fail "Docker daemon is not reachable. Start Docker Desktop or check Docker permissions, then retry."

if docker compose version >/dev/null 2>&1; then
    COMPOSE=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose)
else
    fail "Docker Compose is required. Install docker compose or docker-compose."
fi

if [[ ! -f server/config.docker.yaml ]]; then
    log "server/config.docker.yaml not found; creating it from server/config-temp.yaml."
    cp server/config-temp.yaml server/config.docker.yaml
fi

if [[ "$PULL" == true ]]; then
    log "Pulling service images..."
    "${COMPOSE[@]}" pull mysql || warn "Pull failed; continuing with local images."
fi

services=(mysql server web)
if [[ "$BUILD" == true ]]; then
    log "Building and starting Docker services: ${services[*]}"
    "${COMPOSE[@]}" up -d --build "${services[@]}"
else
    log "Starting Docker services without rebuilding: ${services[*]}"
    "${COMPOSE[@]}" up -d "${services[@]}"
fi

log "Waiting for MySQL healthcheck..."
for _ in {1..60}; do
    mysql_health="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}unknown{{end}}' gshark-mysql 2>/dev/null || true)"
    if [[ "$mysql_health" == "healthy" ]]; then
        break
    fi
    sleep 2
done

mysql_health="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}unknown{{end}}' gshark-mysql 2>/dev/null || true)"
if [[ "$mysql_health" != "healthy" ]]; then
    warn "MySQL is not healthy yet. Check logs with: ${COMPOSE[*]} logs mysql"
fi

if [[ "$WITH_SCAN" == true ]]; then
    log "Starting scan container..."
    if [[ "$BUILD" == true ]]; then
        "${COMPOSE[@]}" up -d --build scan
    else
        "${COMPOSE[@]}" up -d scan
    fi
else
    warn "Scan container was not started. Start it later with: scripts/quick-docker.sh --with-scan --no-build"
fi

log "Current container status:"
"${COMPOSE[@]}" ps

cat <<'EOF'

GShark Docker quick deploy finished.

Open the web UI:
  http://localhost:8080

Default login:
  gshark / gshark

Useful commands:
  docker compose logs -f server
  docker compose restart scan
  docker compose down
EOF
