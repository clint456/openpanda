#!/bin/bash
# ============================================================
# deploy.sh — 服务器运维快捷脚本
#
# 用法:
#   chmod +x deploy.sh
#   ./deploy.sh backup    # 备份数据库
#   ./deploy.sh logs      # 查看服务日志
#   ./deploy.sh status    # 查看服务状态
# ============================================================
set -e

PROJECT_DIR="/data/openpanda"
COMPOSE_FILE="docker-compose.prod.yml"

# 优先使用 docker compose (插件), 其次 docker-compose (独立二进制)
if docker compose version &>/dev/null; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose &>/dev/null && docker-compose version &>/dev/null; then
  COMPOSE_CMD="docker-compose"
else
  echo "❌ 错误: 未找到可用的 docker compose"
  exit 1
fi

cd "$PROJECT_DIR"

backup_db() {
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="/data/openpanda/backups/openpanda_${TIMESTAMP}.sql"
    echo "========== 备份数据库到: $BACKUP_FILE =========="
    docker exec openpanda-db pg_dump -U postgres openpanda > "$BACKUP_FILE"
    echo "✓ 备份完成 ($(du -h "$BACKUP_FILE" | cut -f1))"
}

show_logs() {
    $COMPOSE_CMD -f "$COMPOSE_FILE" logs -f --tail=100
}

status() {
    $COMPOSE_CMD -f "$COMPOSE_FILE" ps
    echo ""
    echo "磁盘使用:"
    df -h /data
}

case "${1}" in
    backup) backup_db ;;
    logs)   show_logs ;;
    status) status ;;
    *)      echo "用法: $0 {backup|logs|status}" ;;
esac
