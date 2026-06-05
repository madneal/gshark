#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO="${GSHARK_REPO:-madneal/gshark}"
API_BASE="${GSHARK_API_BASE:-https://api.github.com/repos/$REPO/releases}"
WORK_DIR="${GSHARK_WORK_DIR:-/tmp/gshark-quick-deploy}"
FRONTEND_PORT="${GSHARK_FRONTEND_PORT:-8080}"
BACKEND_PORT="${GSHARK_BACKEND_PORT:-8888}"
LOCAL_FILE=""
ASSET_URL=""

usage() {
    cat <<'EOF'
Usage: scripts/quick-release.sh [options]

Quickly deploy GShark from a release zip. This configures Nginx for the web UI
and starts the bundled gshark backend in the background.

Options:
  --file PATH         Use a local release zip instead of downloading one.
  --asset-url URL    Download a specific release asset URL.
  --work-dir PATH    Directory used for temporary release files.
  -h, --help         Show this help message.

Environment:
  GSHARK_FRONTEND_PORT  Nginx listen port. Default: 8080
  GSHARK_BACKEND_PORT   Backend port. Default: 8888
  GSHARK_HTML_ROOT      Nginx web root. Default: OS-specific

Examples:
  scripts/quick-release.sh
  scripts/quick-release.sh --file ./gshark_linux_amd64.zip
EOF
}

log() {
    printf '[INFO] %s\n' "$*" >&2
}

warn() {
    printf '[WARN] %s\n' "$*" >&2
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

require_cmd() {
    command -v "$1" >/dev/null 2>&1 || fail "$1 is required."
}

run_maybe_sudo() {
    if "$@"; then
        return 0
    fi

    if sudo -n true >/dev/null 2>&1; then
        sudo "$@"
        return
    fi

    if [[ -t 0 ]]; then
        sudo "$@"
        return
    fi

    fail "Administrator privileges are required to run: $*"
}

require_writable_or_sudo() {
    local path="$1"
    local action="$2"

    if [[ -w "$path" ]]; then
        return 0
    fi

    if sudo -n true >/dev/null 2>&1; then
        return 0
    fi

    if [[ -t 0 ]]; then
        log "$action requires administrator privileges."
        sudo -v
        return
    fi

    fail "$action requires administrator privileges for $path. Re-run from an interactive terminal or grant sudo first."
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --file)
            [[ $# -ge 2 ]] || fail "--file requires a path."
            LOCAL_FILE="$2"
            shift
            ;;
        --asset-url)
            [[ $# -ge 2 ]] || fail "--asset-url requires a URL."
            ASSET_URL="$2"
            shift
            ;;
        --work-dir)
            [[ $# -ge 2 ]] || fail "--work-dir requires a path."
            WORK_DIR="$2"
            shift
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

detect_platform() {
    local os arch
    case "$(uname -s)" in
        Darwin) os="darwin" ;;
        Linux) os="linux" ;;
        *) fail "Unsupported OS: $(uname -s). Use --file with a compatible release package." ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) arch="$(uname -m)" ;;
    esac

    printf '%s_%s' "$os" "$arch"
}

select_release_asset() {
    require_cmd curl
    require_cmd jq

    local platform api_response names urls exact_index fallback_index
    platform="$(detect_platform)"
    log "Fetching latest release metadata from $REPO..."
    api_response="$(curl -fsSL "$API_BASE/latest")" || fail "Failed to fetch latest release metadata."

    names=()
    while IFS= read -r name; do
        names+=("$name")
    done < <(printf '%s' "$api_response" | jq -r '.assets[].name')

    urls=()
    while IFS= read -r url; do
        urls+=("$url")
    done < <(printf '%s' "$api_response" | jq -r '.assets[].browser_download_url')
    [[ ${#names[@]} -gt 0 ]] || fail "No release assets found."

    exact_index=""
    fallback_index=""
    for i in "${!names[@]}"; do
        if [[ "${names[$i]}" == *"$platform"* && "${names[$i]}" == *.zip ]]; then
            exact_index="$i"
            break
        fi
        if [[ -z "$fallback_index" && "${names[$i]}" == *.zip ]]; then
            fallback_index="$i"
        fi
    done

    if [[ -n "$exact_index" ]]; then
        printf '%s|%s\n' "${names[$exact_index]}" "${urls[$exact_index]}"
        return
    fi

    if [[ -t 0 ]]; then
        warn "No exact release asset matched platform: $platform"
        printf 'Available assets:\n'
        for i in "${!names[@]}"; do
            printf '  %d) %s\n' "$((i + 1))" "${names[$i]}"
        done

        local choice
        while true; do
            read -r -p "Choose a release asset (1-${#names[@]}): " choice
            if [[ "$choice" =~ ^[0-9]+$ ]] && (( choice >= 1 && choice <= ${#names[@]} )); then
                printf '%s|%s\n' "${names[$((choice - 1))]}" "${urls[$((choice - 1))]}"
                return
            fi
            warn "Invalid choice."
        done
    fi

    [[ -n "$fallback_index" ]] || fail "No zip release asset found."
    warn "No exact release asset matched platform: $platform; using ${names[$fallback_index]}"
    printf '%s|%s\n' "${names[$fallback_index]}" "${urls[$fallback_index]}"
}

default_html_root() {
    if [[ "$(uname -s)" == "Darwin" ]]; then
        if command -v brew >/dev/null 2>&1; then
            local brew_prefix
            brew_prefix="$(brew --prefix 2>/dev/null || true)"
            if [[ -n "$brew_prefix" ]]; then
                printf '%s/var/www\n' "$brew_prefix"
                return
            fi
        fi
        printf '/usr/local/var/www\n'
    else
        printf '/var/www/html\n'
    fi
}

write_nginx_config() {
    require_cmd nginx

    local html_root nginx_conf nginx_test backup_file nginx_user tmp_nginx_conf
    html_root="${GSHARK_HTML_ROOT:-$(default_html_root)}"
    nginx_test="$(nginx -t 2>&1 || true)"
    nginx_conf="$(printf '%s\n' "$nginx_test" | awk '/configuration file .* test is successful/ {print $4; exit}')"
    [[ -n "$nginx_conf" ]] || fail "Could not detect nginx.conf. Output was: $nginx_test"

    if [[ "$(uname -s)" == "Darwin" ]]; then
        nginx_user="$(id -un) staff"
    else
        nginx_user="www-data www-data"
    fi

    require_writable_or_sudo "$nginx_conf" "Writing Nginx config"

    backup_file="${nginx_conf}.backup.$(date +%Y%m%d_%H%M%S)"
    tmp_nginx_conf="$WORK_DIR/nginx.conf.quick"
    mkdir -p "$WORK_DIR"

    log "Backing up Nginx config to $backup_file"
    run_maybe_sudo cp "$nginx_conf" "$backup_file"

    cat > "$tmp_nginx_conf" <<EOF
user  $nginx_user;
worker_processes  1;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;

    server {
        listen       $FRONTEND_PORT;
        server_name  localhost;

        location / {
            root   $html_root;
            index  index.html index.htm;
            try_files \$uri \$uri/ /index.html;
        }

        location /api/ {
            rewrite ^/api/(.*)\$ /\$1 break;
            proxy_pass http://127.0.0.1:$BACKEND_PORT;
            proxy_set_header Host \$http_host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }
    }
}
EOF
    log "Writing quick deploy Nginx config to $nginx_conf"
    run_maybe_sudo cp "$tmp_nginx_conf" "$nginx_conf"

    if ! nginx -t; then
        warn "Nginx config test failed. Restoring backup."
        run_maybe_sudo mv "$backup_file" "$nginx_conf"
        fail "Nginx config test failed."
    fi

    printf '%s\n' "$html_root"
}

restart_nginx() {
    if [[ "$(uname -s)" == "Darwin" ]] && command -v brew >/dev/null 2>&1; then
        log "Restarting Homebrew Nginx..."
        brew services restart nginx >/dev/null 2>&1 || {
            warn "brew services restart failed; trying nginx -s reload."
            run_maybe_sudo nginx -s reload || run_maybe_sudo nginx
        }
    elif command -v systemctl >/dev/null 2>&1; then
        log "Restarting systemd Nginx..."
        run_maybe_sudo systemctl restart nginx
    else
        log "Reloading Nginx..."
        run_maybe_sudo nginx -s reload || run_maybe_sudo nginx
    fi
}

stop_process_on_backend_port() {
    if ! command -v lsof >/dev/null 2>&1; then
        warn "lsof is not available; skipping backend port cleanup."
        return
    fi

    local pid
    pid="$(lsof -t -i:"$BACKEND_PORT" || true)"
    if [[ -n "$pid" ]]; then
        warn "Stopping process on port $BACKEND_PORT: $pid"
        kill $pid || true
        sleep 1
    fi
}

download_or_copy_release() {
    require_cmd unzip
    mkdir -p "$WORK_DIR"
    rm -rf "$WORK_DIR"/gshark_*

    local zip_path asset_name selection
    if [[ -n "$LOCAL_FILE" ]]; then
        [[ -f "$LOCAL_FILE" ]] || fail "Local release zip not found: $LOCAL_FILE"
        zip_path="$WORK_DIR/$(basename "$LOCAL_FILE")"
        cp "$LOCAL_FILE" "$zip_path"
    else
        if [[ -z "$ASSET_URL" ]]; then
            selection="$(select_release_asset)"
            asset_name="${selection%%|*}"
            ASSET_URL="${selection#*|}"
        else
            asset_name="$(basename "${ASSET_URL%%\?*}")"
        fi
        zip_path="$WORK_DIR/$asset_name"
        log "Downloading $asset_name..."
        curl -fL -o "$zip_path" "$ASSET_URL"
    fi

    log "Extracting release package..."
    unzip -o -q "$zip_path" -d "$WORK_DIR"

    local app_dir binary_path
    binary_path="$(find "$WORK_DIR" -maxdepth 2 -type f -name gshark -print -quit)"
    [[ -n "$binary_path" ]] || fail "Could not find gshark binary after extracting $zip_path"
    app_dir="$(dirname "$binary_path")"
    [[ -d "$app_dir" ]] || fail "Could not find extracted app directory for $zip_path"
    printf '%s\n' "$app_dir"
}

deploy_frontend() {
    local app_dir="$1"
    local html_root="$2"
    [[ -d "$app_dir/dist" ]] || fail "Release package does not contain dist/: $app_dir"

    log "Deploying frontend to $html_root"
    run_maybe_sudo mkdir -p "$html_root"
    if command -v rsync >/dev/null 2>&1; then
        run_maybe_sudo rsync -a --delete "$app_dir/dist/" "$html_root/"
    else
        run_maybe_sudo cp -R "$app_dir/dist/." "$html_root/"
    fi

    if [[ ! -w "$html_root" ]]; then
        if [[ "$(uname -s)" == "Darwin" ]]; then
            run_maybe_sudo chown -R "$(id -un):staff" "$html_root"
        else
            run_maybe_sudo chown -R www-data:www-data "$html_root" || true
        fi
    fi
}

start_backend() {
    local app_dir="$1"

    if [[ ! -f "$app_dir/config.yaml" && -f "$app_dir/config-temp.yaml" ]]; then
        log "Creating config.yaml from config-temp.yaml"
        cp "$app_dir/config-temp.yaml" "$app_dir/config.yaml"
    fi

    chmod +x "$app_dir/gshark"
    stop_process_on_backend_port

    log "Starting backend on port $BACKEND_PORT..."
    (
        cd "$app_dir"
        ./gshark serve > gshark.log 2>&1 &
        printf '%s\n' "$!" > gshark.pid
    )

    sleep 2
    local pid
    pid="$(cat "$app_dir/gshark.pid")"
    if ! kill -0 "$pid" >/dev/null 2>&1; then
        fail "Backend failed to start. Check log: $app_dir/gshark.log"
    fi

    printf '%s\n' "$pid"
}

main() {
    log "Starting GShark release quick deploy..."
    local app_dir html_root pid
    app_dir="$(download_or_copy_release)"
    html_root="$(write_nginx_config)"
    deploy_frontend "$app_dir" "$html_root"
    restart_nginx
    pid="$(start_backend "$app_dir")"

    cat <<EOF

GShark release quick deploy finished.

Open the web UI:
  http://localhost:$FRONTEND_PORT

Default login:
  gshark / gshark

Backend:
  PID: $pid
  Log: $app_dir/gshark.log
  Stop: kill $pid
EOF
}

main
