#!/bin/sh

exec 2>&1
set -e

print_banner() {
  msg="$1"
  side=30
  content=" $msg "
  content_len=${#content}
  line_len=$((side * 2 + content_len))

  line=$(printf '*%.0s' $(seq 1 "$line_len"))
  side_eq=$(printf '*%.0s' $(seq 1 "$side"))

  printf "%s\n%s%s%s\n%s\n" "$line" "$side_eq" "$content" "$side_eq" "$line"
}

print_banner_delay() {
  msg="$1"
  delay="$2"

  (
    sleep "$delay"
    print_banner "$msg"
  ) &
}

export ROCKETMQ_GO_LOG_LEVEL=error

printf "+ Waiting for basic services - redis, mysql, minio, clickhouse, rocketmq (namesrv & broker) - to stabilize...\n"
sleep 25

if [ "$RUN_MODE" = "debug" ]; then
  print_banner "Cozeloop Server: Starting in [DEBUG] mode..."
  print_banner_delay "Cozeloop Server: Successfully Started in [DEBUG] mode! Please toggle debugger in IDEA at [HOST_IP:40000]." 5

  set -x
  dlv exec /cozeloop-bin/backend/debug/main \
    --headless \
    --listen=:40000 \
    --api-version=2 \
    --accept-multiclient \
    --log

  wait
elif [ "$RUN_MODE" = "release" ]; then
  print_banner "Cozeloop Server: Starting in [RELEASE] mode..."
  print_banner_delay "Cozeloop Server: Successfully Started in [RELEASE] mode!" 5

  set -x
  /cozeloop-bin/backend/release/main

  wait
else
  print_banner "Cozeloop Server: Starting in [DEV] mode..."
  print_banner_delay "Cozeloop Server: Successfully Started in [DEV] mode!" 65

  set -x

  air -c /etc/cozeloop-server/.air.toml

  wait
fi