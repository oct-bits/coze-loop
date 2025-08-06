#!/bin/bash

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

export MINIO_ROOT_USER="${COZELOOP_MINIO_USER}"
export MINIO_ROOT_PASSWORD="${COZELOOP_MINIO_PASSWORD}"
export MC_HOST_myminio="http://${COZELOOP_MINIO_USER}:${COZELOOP_MINIO_PASSWORD}@${COZELOOP_MINIO_APP_NAME}:${COZELOOP_MINIO_PORT_POD}"

print_banner "Minio: Starting..."
print_banner_delay "Minio: Successfully Started!" 5

echo "+ minio server"
minio server /minio_data --address ":${COZELOOP_MINIO_PORT_POD}" &

sleep 5

set -x
if ! mc ls myminio/"${COZELOOP_MINIO_BUCKET}" > /dev/null 2>&1; then
  mc mb myminio/"${COZELOOP_MINIO_BUCKET}"
fi

wait