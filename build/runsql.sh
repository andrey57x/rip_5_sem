#!/usr/bin/env bash
set -euo pipefail

# run_fill.sh - выполнить указанный SQL-файл в Postgres, используя .env при наличии
# Ожидается, что .env (в той же папке что и скрипт) содержит:
# DB_HOST, DB_NAME, DB_PORT, DB_USER, DB_PASS
#
# Usage:
#   ./run_fill.sh /path/to/file.sql

# --- load .env if present (simple sourcing) ---
ENV_FILE="$(dirname "$0")/.env"
if [ -f "$ENV_FILE" ]; then
  # Подключаем .env в текущую среду (только строки вида VAR=val — пользовательский .env должен быть доверенным)
  # set -a временно экспортирует все переменные, затем source и set +a
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

# --- defaults (если переменные не заданы в окружении/.env) ---
DB_HOST="${DB_HOST:-localhost}"
DB_NAME="${DB_NAME:-mass_calculation_db}"
DB_PORT="${DB_PORT:-5433}"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-1337}"

# --- аргументы ---
if [ $# -lt 1 ]; then
  echo "ERROR: SQL filename required."
  echo "Usage: $0 /path/to/file.sql"
  exit 1
fi

SQL_FILE="$1"

# --- проверки ---
if ! command -v psql >/dev/null 2>&1; then
  echo "ERROR: psql not found in PATH."
  exit 2
fi

if [ ! -f "$SQL_FILE" ]; then
  echo "ERROR: SQL file '$SQL_FILE' not found."
  exit 3
fi

# --- безопасность: экспортируем PGPASSWORD только на время выполнения ---
WE_SET_PGPASSWORD=0
cleanup() {
  local rc=$?
  if [ "${WE_SET_PGPASSWORD:-0}" -eq 1 ]; then
    unset PGPASSWORD
  fi
  exit $rc
}
trap cleanup EXIT

if [ -n "${DB_PASS:-}" ]; then
  export PGPASSWORD="$DB_PASS"
  WE_SET_PGPASSWORD=1
fi

# --- выполнить ---
echo "Running '$SQL_FILE' on ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME} ..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$SQL_FILE"

echo "Done."
