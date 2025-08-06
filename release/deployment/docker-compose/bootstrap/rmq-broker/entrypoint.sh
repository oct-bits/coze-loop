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

rmq_home() {
  base_dir="/home/rocketmq"
  for d in "${base_dir}"/rocketmq-*; do
    [ -d "$d" ] && echo "$d" && return
  done
}

print_banner "Starting..."

mkdir -p /store/logs

(
  while true; do
    sleep 3
    if sh /coze-loop-rmq-broker/bootstrap/healthcheck.sh; then
      print_banner "Completed!"
      break
    else
      sleep 1
    fi
  done
)&

exec "$(rmq_home)"/bin/mqbroker -n coze-loop-rmq-namesrv:9876