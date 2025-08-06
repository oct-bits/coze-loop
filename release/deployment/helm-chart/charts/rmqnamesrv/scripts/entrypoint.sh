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

rmq_home() {
  base_dir="/home/rocketmq"
  for d in "$base_dir"/rocketmq-*; do
    [ -d "$d" ] && echo "$d" && return
  done
}

print_banner "RmqNamesrv: Starting..."
print_banner_delay "RmqNamesrv: Successfully Started!" 3

echo "+ mkdir -p /store/logs"
mkdir -p /store/logs

echo "+ mqnamesrv"
exec "$(rmq_home)/bin/mqnamesrv"