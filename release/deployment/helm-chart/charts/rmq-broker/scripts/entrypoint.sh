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

# 异步延迟打印 Banner
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

ROCKETMQ_HOME="$(rmq_home)"
MQBROKER_CMD="${ROCKETMQ_HOME}/bin/mqbroker"
MQADMIN_CMD="${ROCKETMQ_HOME}/bin/mqadmin"
MQNAMESRV_ADDR="cozeloop-rmqnamesrv:9876"  # 用unbrella统一传进来

declare -A topics
{
  while IFS='=' read -r topic consumers || [[ -n "$topic" ]]; do
    [[ -z "$topic" || "${topic:0:1}" == "#" ]] && continue
    topics["$topic"]="$consumers"
  done
} < /etc/cozeloop-rmqbroker/topics.cfg

print_banner "RmqBroker: Starting..."
print_banner_delay "RmqBroker: Successfully Started!" 35

echo "+ mkdir -p /store/logs"
mkdir -p /store/logs

echo "+ mqbroker"
"$MQBROKER_CMD" -n "${MQNAMESRV_ADDR}" &

sleep 10

i=1
for topic in "${!topics[@]}"; do
  ii=$i
  (
    echo "+ Check if topic#$ii('$topic') exists...: mqadmin topicList | grep -q '^$topic$'"
    if ! "$MQADMIN_CMD" topicList -n "$MQNAMESRV_ADDR" | grep -q "^$topic$"; then
      echo "[+] Topic#$ii('$topic') not exists, now creating...: mqadmin updateTopic -t $topic -r 8 -w 8"
      "$MQADMIN_CMD" updateTopic -n "$MQNAMESRV_ADDR" -c DefaultCluster -t "$topic" -r 8 -w 8
    else
      echo "[-] Topic#$ii('$topic') already exists."
    fi

    IFS=',' read -ra consumer_groups <<< "${topics[$topic]}"
    j=1
    for group in "${consumer_groups[@]}"; do
      echo "++ Check if consumer#$ii-$j('$group') exists...: mqadmin consumerProgress | grep -q '^$group$'"
      if ! "$MQADMIN_CMD" consumerProgress -n "$MQNAMESRV_ADDR" | grep -q "^$group$"; then
        echo "[++] Consumer#$ii-$j('$group') not exists, now creating...: mqadmin updateSubGroup -g $group"
        "$MQADMIN_CMD" updateSubGroup -n "$MQNAMESRV_ADDR" -c DefaultCluster -g "$group"

        retry_topic="%RETRY%$group"
        echo "[+++] Consumer#$ii-$j('$group')'s related retry topic('$retry_topic') is creating...: mqadmin updateTopic -t $retry_topic -r 8 -w 8"
        "$MQADMIN_CMD" updateTopic -n "$MQNAMESRV_ADDR" -c DefaultCluster -t "$retry_topic" -r 8 -w 8
      else
        echo "[--] Consumer#$ii-$j('$group')' already exists."
      fi
      j=$((j + 1))
    done

    echo "+ Topic#$ii('$topic') is ready! (with it's consumers and retry topics)"
  ) &
  i=$((i + 1))
done

echo "+ All topics have been send in batch"

wait