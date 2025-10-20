#!/bin/sh
# 此脚本会在 nginx 启动前被 /docker-entrypoint.sh 执行
# 它在后台启动一个监控进程，监控证书文件变化并自动重载 nginx

# 安装 inotify-tools
if ! command -v inotifywait > /dev/null; then
    echo "Installing inotify-tools..."
    apk add --no-cache inotify-tools
fi

# 获取证书目录
CERT_DIR="/etc/letsencrypt/live/${SECOND_LEVEL_DOMAIN_NAME}"
CERT_FILE="${CERT_DIR}/fullchain.pem"

echo "Starting certificate file monitor for: ${CERT_FILE}"

# 后台监控证书文件变化
(
    # 等待 nginx 启动
    sleep 10

    while true; do
        if [ -f "$CERT_FILE" ]; then
            # 监控证书文件和目录(证书更新时可能是替换整个目录)
            inotifywait -e modify,create,moved_to,attrib "$CERT_FILE" "$CERT_DIR" 2>/dev/null

            if [ $? -eq 0 ]; then
                echo "$(date): Certificate updated, reloading nginx..."
                sleep 2
                nginx -s reload && echo "$(date): Nginx reloaded successfully" || echo "$(date): Nginx reload failed"
            fi
        else
            echo "$(date): Certificate file not found, waiting..."
            sleep 60
        fi
    done
) &

echo "Certificate monitor started in background"
