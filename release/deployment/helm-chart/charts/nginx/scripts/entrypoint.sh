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

print_banner "Cozeloop Ingress: Starting..."
print_banner_delay "Cozeloop Ingress: Successfully Started!" 3

if [ -n "$COZELOOP_OSS_ADDR_POD_CUSTOMIZED" ]; then
  export COZELOOP_OSS_ADDR_POD="$COZELOOP_OSS_ADDR_POD_CUSTOMIZED"
fi
if [ -n "$COZELOOP_OSS_BUCKET_CUSTOMIZED" ]; then
  export COZELOOP_OSS_BUCKET="$COZELOOP_OSS_BUCKET_CUSTOMIZED"
fi

set -x

rm -f /etc/nginx/nginx.conf
mkdir -p /etc/nginx
sed -e "s|{{\$COZELOOP_APP_HERTZ_ADDR_POD}}|${COZELOOP_APP_HERTZ_ADDR_POD}|g" \
    -e "s|{{\$COZELOOP_OSS_ADDR_POD}}|${COZELOOP_OSS_ADDR_POD}|g" \
    -e "s|{{\$COZELOOP_OSS_BUCKET}}|${COZELOOP_OSS_BUCKET}|g" \
    /etc/cozeloop-ingress/server.conf > /etc/nginx/nginx.conf
chmod 444 /etc/nginx/nginx.conf

exec /docker-entrypoint.sh nginx -g 'daemon off;'