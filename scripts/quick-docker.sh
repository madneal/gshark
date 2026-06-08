#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

WITH_SCAN=false

usage() {
    cat <<'EOF'
Usage: scripts/quick-docker.sh [--with-scan]

Build and start GShark with Docker Compose.

Options:
  --with-scan  Also start the scanner container.
  -h, --help   Show this help message.
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --with-scan)
            WITH_SCAN=true
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

echo
"${COMPOSE[@]}" ps
echo
echo "GShark is starting at: http://localhost:8080"
echo "Default login: gshark / gshark"
