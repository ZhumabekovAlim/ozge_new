#!/bin/sh
set -e
HOST="${DB_HOST:-db}"
PORT="${DB_PORT:-3306}"

while ! nc -z "$HOST" "$PORT" >/dev/null 2>&1; do
  echo "Waiting for MySQL at $HOST:$PORT..."
  sleep 1
done

exec "$@"
