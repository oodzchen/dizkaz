#!/bin/bash
set -e

# 默认环境文件
ENV_FILE=".env.local"
COMMAND=""
STEPS=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--env)
            ENV_FILE="$2"
            shift 2
            ;;
        up|down|goto|force|version|status)
            COMMAND="$1"
            shift
            if [[ $# -gt 0 && $1 =~ ^[0-9]+$ ]]; then
                STEPS="$1"
                shift
            fi
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS] COMMAND [STEPS]"
            echo ""
            echo "Options:"
            echo "  -e, --env FILE    Environment file (default: .env.local)"
            echo "  -h, --help        Show this help message"
            echo ""
            echo "Commands:"
            echo "  up [N]           Apply all or N pending migrations"
            echo "  down [N]         Roll back all or N migrations"
            echo "  goto V           Migrate to specific version"
            echo "  force V          Set version V but don't run migration"
            echo "  version          Show current migration version"
            echo "  status           Alias for version command"
            echo ""
            echo "Examples:"
            echo "  $0 up                    # Apply all pending migrations"
            echo "  $0 -e .env.prod up 1    # Apply 1 migration using production env"
            echo "  $0 down 1               # Roll back 1 migration"
            echo "  $0 version              # Show current migration version"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use $0 --help for usage information"
            exit 1
            ;;
    esac
done

# 检查命令是否提供
if [[ -z "$COMMAND" ]]; then
    echo "Error: No command specified"
    echo "Use $0 --help for usage information"
    exit 1
fi

# 检查环境文件是否存在
if [[ ! -f "$ENV_FILE" ]]; then
    echo "Error: Environment file '$ENV_FILE' not found"
    exit 1
fi

# 加载环境变量
source ./scripts/load-env.sh
readenv "$ENV_FILE"

# 检查 migrate 工具是否安装
if ! command -v migrate &> /dev/null; then
    echo "Error: migrate tool is not installed"
    echo "Install with: brew install golang-migrate"
    exit 1
fi

# 构建数据库连接字符串
DB_CONNECTION="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# 执行命令
echo "Running migrate $COMMAND $STEPS with environment: $ENV_FILE"
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo ""

case "$COMMAND" in
    up|down)
        if [[ -n "$STEPS" ]]; then
            migrate -path config/db/migrations -database "$DB_CONNECTION" "$COMMAND" "$STEPS"
        else
            migrate -path config/db/migrations -database "$DB_CONNECTION" "$COMMAND"
        fi
        ;;
    goto|force)
        if [[ -z "$STEPS" ]]; then
            echo "Error: $COMMAND command requires a version number"
            exit 1
        fi
        migrate -path config/db/migrations -database "$DB_CONNECTION" "$COMMAND" "$STEPS"
        ;;
    version)
        migrate -path config/db/migrations -database "$DB_CONNECTION" "$COMMAND"
        ;;
    status)
        migrate -path config/db/migrations -database "$DB_CONNECTION" version
        ;;
    *)
        echo "Unknown command: $COMMAND"
        exit 1
        ;;
esac

echo ""
echo "Migration $COMMAND completed successfully!"