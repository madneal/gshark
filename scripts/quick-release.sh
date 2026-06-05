#!/usr/bin/env bash
set -euo pipefail

REPO="${GSHARK_REPO:-madneal/gshark}"
WORK_DIR="${GSHARK_WORK_DIR:-/tmp/gshark-quick-release}"
FRONTEND_PORT="${GSHARK_FRONTEND_PORT:-8080}"
BACKEND_PORT="${GSHARK_BACKEND_PORT:-8888}"
ZIP_FILE=""

usage() {
    cat <<'EOF'
Usage: scripts/quick-release.sh [--file PATH]

Download a GShark release package, deploy dist/ to Nginx, and start the backend.

Options:
  --file PATH  Use a local release zip instead of downloading the latest release.
  -h, --help   Show this help message.

Environment:
  GSHARK_FRONTEND_PORT  Frontend port. Default: 8080
  GSHARK_BACKEND_PORT   Backend port. Default: 8888
  GSHARK_HTML_ROOT      Nginx web root. Default: Homebrew var/www on macOS,
                        /var/www/html on Linux.
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --file)
            [[ $# -ge 2 ]] || {
                echo "--file requires a path." >&2
                exit 1
            }
            ZIP_FILE="$2"
            shift
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

require() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "$1 is required." >&2
        exit 1
    }
}

platform() {
    case "$(uname -s)" in
        Darwin) os="darwin" ;;
        Linux) os="linux" ;;
        *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) echo "Unsupported arch: $(uname -m)" >&2; exit 1 ;;
    esac

    echo "${os}_${arch}"
}

html_root() {
    if [[ -n "${GSHARK_HTML_ROOT:-}" ]]; then
        echo "$GSHARK_HTML_ROOT"
    elif [[ "$(uname -s)" == "Darwin" ]] && command -v brew >/dev/null 2>&1; then
        echo "$(brew --prefix)/var/www"
    elif [[ "$(uname -s)" == "Darwin" ]]; then
        echo "/usr/local/var/www"
    else
        echo "/var/www/html"
    fi
}

require curl
require jq
require unzip
require nginx

mkdir -p "$WORK_DIR"
rm -rf "$WORK_DIR"/gshark_*

if [[ -z "$ZIP_FILE" ]]; then
    target="$(platform)"
    echo "[INFO] Downloading latest ${target} release..."
    asset_line="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | jq -r --arg target "$target" '.assets[] | select(.name | contains($target) and endswith(".zip")) | [.name, .browser_download_url] | @tsv' | head -n 1)"
    [[ -n "$asset_line" ]] || {
        echo "No release zip matched platform: $target" >&2
        exit 1
    }

    IFS=$'\t' read -r asset_name asset_url <<< "$asset_line"
    ZIP_FILE="$WORK_DIR/$asset_name"
    curl -fL -o "$ZIP_FILE" "$asset_url"
fi

[[ -f "$ZIP_FILE" ]] || {
    echo "Release zip not found: $ZIP_FILE" >&2
    exit 1
}

echo "[INFO] Extracting $ZIP_FILE..."
unzip -o -q "$ZIP_FILE" -d "$WORK_DIR"
BINARY_PATH="$(find "$WORK_DIR" -maxdepth 2 -type f -name gshark -print -quit)"
[[ -n "$BINARY_PATH" ]] || {
    echo "Could not find gshark binary in release package." >&2
    exit 1
}
APP_DIR="$(dirname "$BINARY_PATH")"

HTML_ROOT="$(html_root)"
NGINX_CONF="$(nginx -t 2>&1 | awk '/configuration file .* test is successful/ {print $4; exit}')"
[[ -n "$NGINX_CONF" ]] || {
    echo "Could not locate nginx.conf from nginx -t output." >&2
    exit 1
}

if [[ "$(uname -s)" == "Darwin" ]]; then
    NGINX_USER="$(id -un) staff"
else
    NGINX_USER="www-data www-data"
fi

echo "[INFO] Writing Nginx config: $NGINX_CONF"
sudo cp "$NGINX_CONF" "${NGINX_CONF}.backup.$(date +%Y%m%d_%H%M%S)"
sudo tee "$NGINX_CONF" >/dev/null <<EOF
user  $NGINX_USER;
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
            root   $HTML_ROOT;
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

sudo nginx -t

echo "[INFO] Deploying frontend to $HTML_ROOT..."
sudo mkdir -p "$HTML_ROOT"
sudo cp -R "$APP_DIR/dist/." "$HTML_ROOT/"

if [[ "$(uname -s)" == "Darwin" ]] && command -v brew >/dev/null 2>&1; then
    brew services restart nginx
elif command -v systemctl >/dev/null 2>&1; then
    sudo systemctl restart nginx
else
    sudo nginx -s reload || sudo nginx
fi

if [[ -f "$APP_DIR/config-temp.yaml" && ! -f "$APP_DIR/config.yaml" ]]; then
    cp "$APP_DIR/config-temp.yaml" "$APP_DIR/config.yaml"
fi

if command -v lsof >/dev/null 2>&1; then
    old_pid="$(lsof -t -i:"$BACKEND_PORT" || true)"
    [[ -z "$old_pid" ]] || kill $old_pid
fi

chmod +x "$APP_DIR/gshark"
(
    cd "$APP_DIR"
    ./gshark serve > gshark.log 2>&1 &
    echo "$!" > gshark.pid
)

PID="$(cat "$APP_DIR/gshark.pid")"

echo
echo "GShark is starting at: http://localhost:$FRONTEND_PORT"
echo "Default login: gshark / gshark"
echo "Backend PID: $PID"
echo "Backend log: $APP_DIR/gshark.log"
