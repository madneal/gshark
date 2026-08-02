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

Build GShark's server, web, and scanner images, start the server stack,
then initialize the database if needed (admin account via flags — no browser required).

Options:
  --with-scan              Also start the scanner container.
  --admin-user NAME        Admin login username (default: gshark).
  --admin-password PASS    Admin login password (default: gshark).
  --skip-init              Do not run gshark init after start.
  -h, --help               Show this help message.

gshark init exit codes (used by this script):
  0  applied successfully (server is restarted so serve reconnects)
  2  already initialized (credentials not changed)
  1  failure

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

echo "[INFO] Building server/web/scan images..."
"${COMPOSE[@]}" build server web scan

echo "[INFO] Starting mysql..."
"${COMPOSE[@]}" up -d mysql

echo "[INFO] Starting server/web..."
"${COMPOSE[@]}" up -d server web

INIT_RESULT="skipped" # skipped | applied | failed | skipped_flag

if [[ "$SKIP_INIT" == true ]]; then
    INIT_RESULT="skipped_flag"
else
    echo "[INFO] Waiting for MySQL healthy and server binary..."
    ready=false
    for i in $(seq 1 60); do
        h=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' gshark-mysql 2>/dev/null || echo missing)
        if [[ "$h" == "healthy" ]]; then
            if docker exec gshark-server test -x ./gshark 2>/dev/null; then
                ready=true
                break
            fi
        fi
        sleep 2
    done
    if [[ "$ready" != true ]]; then
        echo "[ERROR] Timed out waiting for MySQL healthy / gshark-server binary." >&2
        exit 1
    fi

    echo "[INFO] Initializing database (admin-user=${ADMIN_USER})..."
    # Separate process from long-lived serve: on success we must restart server
    # so GVA_DB reconnects (otherwise NeedInit blocks login).
    set +e
    docker exec gshark-server ./gshark init \
        --host "$MYSQL_HOST" \
        --port "$MYSQL_PORT" \
        --user "$MYSQL_USER" \
        --password "$MYSQL_PASSWORD" \
        --db "$MYSQL_DB" \
        --admin-user "$ADMIN_USER" \
        --admin-password "$ADMIN_PASSWORD"
    init_rc=$?
    set -e

    case "$init_rc" in
        0)
            INIT_RESULT="applied"
            echo "[INFO] Init applied; restarting server so it reconnects to MySQL..."
            "${COMPOSE[@]}" restart server
            # wait for server to accept traffic again
            for i in $(seq 1 30); do
                if curl -sS -o /dev/null -w '' -X POST "http://localhost:8888/init/checkdb" 2>/dev/null; then
                    break
                fi
                sleep 1
            done
            ;;
        2)
            INIT_RESULT="skipped"
            echo "[INFO] Database already initialized; admin credentials were not changed."
            ;;
        *)
            INIT_RESULT="failed"
            echo "[ERROR] gshark init failed (exit ${init_rc})." >&2
            echo "        Open http://localhost:8080 to initialize in the UI, or re-run:" >&2
            echo "        docker exec gshark-server ./gshark init --help" >&2
            ;;
    esac
fi

if [[ "$WITH_SCAN" == true && "$INIT_RESULT" != "failed" ]]; then
    echo "[INFO] Starting scan after database initialization..."
    "${COMPOSE[@]}" up -d scan
fi

echo
"${COMPOSE[@]}" ps
echo
echo "GShark is starting at: http://localhost:8080"
case "$INIT_RESULT" in
    applied)
        echo "Admin login: ${ADMIN_USER} / (password from --admin-password)"
        ;;
    skipped)
        echo "Admin login: unchanged (DB was already initialized; --admin-user/--admin-password ignored)"
        ;;
    skipped_flag)
        echo "Admin login: not initialized by this script (--skip-init)"
        ;;
    failed)
        echo "Admin login: unknown (init failed — do not assume credentials above were applied)"
        exit 1
        ;;
esac
