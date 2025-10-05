#!/usr/bin/env bash

# Exit script if any command fails
set -e

# ==============================================================================
#  gshark 自动化部署与启动脚本 (v12 - Rename Config)
#
#  更新:
#  - 自动将 config-temp.yaml 重命名为 config.yaml，确保后端服务能正确启动。
#
#  依赖: curl, jq, nginx, lsof
# ==============================================================================

# --- Function to download gshark by showing an interactive menu ---
download_gshark() {
    echo "[INFO] 从 GitHub 获取最新的 release 列表..."
    local REPO="madneal/gshark"
    local API_URL="https://api.github.com/repos/$REPO/releases/latest"
    local api_response
    api_response=$(curl -s "$API_URL")
    if [[ -z "$api_response" ]]; then
        echo "[ERROR] 从 GitHub 获取 release 数据失败。"
        exit 1
    fi
    local asset_names=()
    while IFS= read -r line; do
        asset_names+=("$line")
    done < <(echo "$api_response" | jq -r '.assets[].name')
    local asset_urls=()
    while IFS= read -r line; do
        asset_urls+=("$line")
    done < <(echo "$api_response" | jq -r '.assets[].browser_download_url')
    if [[ ${#asset_names[@]} -eq 0 ]]; then
        echo "[ERROR] 未找到任何 release 文件。"
        exit 1
    fi
    echo "[INFO] 请选择要下载的 release 包:"
    local i=0
    for name in "${asset_names[@]}"; do
        echo "  $((i+1))) $name"
        i=$((i+1))
    done
    echo
    local choice
    while true; do
        read -p "请输入你的选择 (1-${#asset_names[@]}): " choice
        if [[ "$choice" =~ ^[0-9]+$ ]] && [[ "$choice" -ge 1 ]] && [[ "$choice" -le ${#asset_names[@]} ]]; then
            break
        else
            echo "[WARN] 输入无效, 请输入 1 到 ${#asset_names[@]} 之间的数字。"
        fi
    done
    local selected_index=$((choice - 1))
    download_url="${asset_urls[$selected_index]}"
    GSHARK_ZIP="${asset_names[$selected_index]}"
    echo "[INFO] 开始下载: '${GSHARK_ZIP}'..."
    curl -L -o "$GSHARK_ZIP" "$download_url"
    if [[ $? -ne 0 ]]; then
        echo "[ERROR] 下载失败。"
        exit 1
    fi
    echo "[INFO] 下载成功。"
}

# --- Main deployment logic ---
main() {
    echo "--- gshark 自动化部署与启动开始 ---"

    # 1. Pre-flight Checks & Cleanup
    echo "[INFO] 执行启动前检查与清理..."
    echo "[INFO] 清理旧的 gshark 压缩包和解压目录..."
    rm -f gshark_*.tar.gz gshark_*.zip
    rm -rf gshark_*-*-*

    if ! command -v lsof &> /dev/null; then
        echo "[WARN] 'lsof' 命令未找到，无法检查端口占用情况。将跳过此步骤。"
    else
        PID_ON_PORT=$(lsof -t -i:8888 || true)
        if [[ -n "$PID_ON_PORT" ]]; then
            echo "[WARN] 发现旧的 gshark 进程 (PID: $PID_ON_PORT) 正在使用端口 8888。"
            echo "[INFO] 正在停止该进程以避免冲突..."
            kill "$PID_ON_PORT"
            sleep 1
            echo "[INFO] 旧进程已停止。"
        else
            echo "[INFO] 端口 8888 未被占用，检查通过。"
        fi
    fi

    # 2. 检查核心依赖
    if ! command -v jq &> /dev/null || ! command -v curl &> /dev/null || ! command -v nginx &> /dev/null; then
        echo "[ERROR] 'curl', 'jq', 'nginx' 是必需的。请先安装它们。"
        exit 1
    fi
    echo "[INFO] 所有核心依赖项检查通过。"

    # 3. 自动下载
    download_gshark

    # 4. 获取 Nginx 路径并确定系统类型
    echo "[INFO] 检测 Nginx 和系统环境..."
    NGINX_CONF=$(nginx -t 2>&1 | grep "test is successful" | awk '{print $4}')
    local NGINX_CONF_USER=""
    local CHOWN_OWNER=""

    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        HTML_ROOT="/usr/local/var/www"
        NGINX_CONF_USER="$(whoami) staff"
        CHOWN_OWNER="$(whoami):staff"
        echo "[INFO] 系统: macOS, Web 根目录: $HTML_ROOT, 配置文件: $NGINX_CONF"
    else
        OS="linux"
        HTML_ROOT="/var/www/html"
        NGINX_CONF_USER="www-data www-data"
        CHOWN_OWNER="www-data:www-data"
        echo "[INFO] 系统: Linux, Web 根目录: $HTML_ROOT, 配置文件: $NGINX_CONF"
    fi

    # 5. 备份并创建 Nginx 配置
    echo "[INFO] 备份当前 Nginx 配置..."
    BACKUP_FILE="${NGINX_CONF}.backup.$(date +%Y%m%d_%H%M%S)"
    sudo cp "$NGINX_CONF" "$BACKUP_FILE"
    echo "[INFO] 备份完成: $BACKUP_FILE"
    echo "[INFO] 创建新的 Nginx 配置..."
    sudo tee "$NGINX_CONF" > /dev/null << EOF
user  $NGINX_CONF_USER;
worker_processes  1;
events { worker_connections  1024; }
http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;
    server {
        listen       8080;
        server_name  localhost;
        location / { root $HTML_ROOT; index index.html index.htm; autoindex on; }
        location /api/ {
            rewrite ^/api/(.*)\$ /\$1 break;
            proxy_pass http://127.0.0.1:8888;
            proxy_set_header Host \$http_host;
            proxy_set_header X-Real-IP \$remote_addr;
        }
    }
}
EOF

    # 6. 测试配置
    echo "[INFO] 测试新的 Nginx 配置..."
    if ! sudo nginx -t; then
        echo "[ERROR] Nginx 配置测试失败。请检查 $NGINX_CONF"
        echo "[INFO] 你的旧配置已恢复。"
        sudo mv "$BACKUP_FILE" "$NGINX_CONF"
        exit 1
    fi

    # 7. 停止并启动 Nginx 服务
    echo "[INFO] 正在应用新的 Nginx 配置..."
    if [[ "$OS" == "macos" ]]; then
        echo "[INFO] 停止 Homebrew Nginx 服务 (如果正在运行)..."
        brew services stop nginx || true
        echo "[INFO] 启动 Homebrew Nginx 服务..."
        brew services start nginx
    else
        echo "[INFO] 停止 systemd Nginx 服务 (如果正在运行)..."
        sudo systemctl stop nginx || true
        echo "[INFO] 启动 systemd Nginx 服务..."
        sudo systemctl start nginx
    fi

    # 8. 部署前端文件
    echo "[INFO] 部署前端文件到 $HTML_ROOT ..."
    unzip -o -q "$GSHARK_ZIP"
    GSHARK_DIR=$(find . -maxdepth 1 -type d -name "gshark*")

    sudo mkdir -p "$HTML_ROOT"
    sudo cp -r "$GSHARK_DIR/dist/"* "$HTML_ROOT/"
    sudo chown -R "$CHOWN_OWNER" "$HTML_ROOT"
    echo "[INFO] 前端文件部署完成。"

    # 9. 配置并启动后端服务
    echo "[INFO] 准备启动后端服务..."
    cd "$GSHARK_DIR"

    # <-- NEW: Rename config file for the service to use -->
    if [[ -f "config-temp.yaml" ]]; then
        echo "[INFO] 正在将 'config-temp.yaml' 重命名为 'config.yaml'..."
        mv config-temp.yaml config.yaml
    else
        echo "[WARN] 未找到 'config-temp.yaml'。将继续启动，服务可能使用默认配置。"
    fi

    echo "[INFO] 赋予 gshark 执行权限..."
    chmod +x gshark

    echo "[INFO] 正在后台启动 gshark 服务 (日志写入 gshark.log)..."
    ./gshark serve > gshark.log 2>&1 &
    GSHARK_PID=$!

    cd ..

    sleep 2

    if ! kill -0 $GSHARK_PID > /dev/null 2>&1; then
        echo "[ERROR] gshark 后端服务启动失败！请检查 '$(pwd)/$GSHARK_DIR/gshark.log' 获取详细信息。"
        exit 1
    fi
    echo "[INFO] gshark 后端服务已成功启动 (PID: $GSHARK_PID)。"

    echo
    echo "--- ✅ 全自动部署与启动成功! ---"
    echo "前端和后端服务均已在后台运行。"
    echo
    echo "  - 访问应用: http://localhost:8080"
    echo "  - 后端 PID: $GSHARK_PID"
    echo "  - 日志文件: $(pwd)/$GSHARK_DIR/gshark.log"
    echo "  - 停止后端: kill $GSHARK_PID"
    echo
}

# 执行主函数
main