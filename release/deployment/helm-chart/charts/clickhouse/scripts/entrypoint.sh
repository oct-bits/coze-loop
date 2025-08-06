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

server_conf() {
  server_conf_path=/etc/clickhouse-server/config.xml
  mkdir -p "$(dirname "$server_conf_path")"
  cat <<EOF > "$server_conf_path"
$(eval "echo \"$(cat "/etc/cozeloop-clickhouse/server.cfg.xml")\"")
EOF
  echo "$server_conf_path"
}

init_db_client_conf() {
  init_db_client_conf_path=/etc/clickhouse-client/init_db.config.xml
  mkdir -p "$(dirname "$init_db_client_conf_path")"
  cat <<EOF > "$init_db_client_conf_path"
$(eval "echo \"$(cat "/etc/cozeloop-clickhouse/init_db.client.cfg.xml")\"")
EOF
  echo "$init_db_client_conf_path"
}

init_table_client_conf() {
  init_table_client_conf_path=/etc/clickhouse-client/init_table.config.xml
  mkdir -p "$(dirname "$init_table_client_conf_path")"
  cat <<EOF > "$init_table_client_conf_path"
$(eval "echo \"$(cat "/etc/cozeloop-clickhouse/init_table.client.cfg.xml")\"")
EOF
  echo "$init_table_client_conf_path"
}

print_banner "Clickhouse: Starting..."
print_banner_delay "Clickhouse: Successfully Started!" 5

server_config_path="$(server_conf)"
echo "+ clickhouse-server"
clickhouse-server --config="$server_config_path" &

sleep 2

init_db_client_config_path="$(init_db_client_conf)"
echo "+ init database: clickhouse-client --query \"CREATE DATABASE IF NOT EXISTS \`${COZELOOP_CLICKHOUSE_DATABASE}\`;\""
clickhouse-client --config "$init_db_client_config_path" --query "CREATE DATABASE IF NOT EXISTS \`${COZELOOP_CLICKHOUSE_DATABASE}\`;"
rm -f "$init_db_client_config_path"

init_table_client_config_path="$(init_table_client_conf)"
i=1
# shellcheck disable=SC2010
for file in $(ls /etc/cozeloop-clickhouse/init-sql | grep '\.sql$'); do
  echo "+ init #$i: clickhouse-client < $file"
  clickhouse-client --config "$init_table_client_config_path" < "/etc/cozeloop-clickhouse/init-sql/$file"
  i=$((i + 1))
done
rm -f "$init_table_client_config_path"

wait