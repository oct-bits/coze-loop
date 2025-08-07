#!/bin/sh

set -e

for i in $(seq 1 60); do
  if redis-cli -h coze-loop-redis -a "${COZE_LOOP_REDIS_PASSWORD}" --no-auth-warning ping | grep -q PONG; then
    echo "[INFO] Redis is ready"
    exit 0
  else
    echo "[INFO] [$i/60] Waiting for Redis..."
    sleep 1
  fi
done

echo "[ERROR] Redis did not become ready after 60 attempts."
exit 1