#!/bin/bash

set -e

init_db_client_conf() {
  init_db_client_conf_path=/etc/clickhouse-client/init_db.config.xml
  mkdir -p "$(dirname "$init_db_client_conf_path")"
  cat <<EOF > "$init_db_client_conf_path"
$(eval "echo \"$(cat "/etc/cozeloop-clickhouse/init_db.client.cfg.xml")\"")
EOF
  echo "$init_db_client_conf_path"
}

init_db_client_config_path="$(init_db_client_conf)"
clickhouse-client --config "$init_db_client_config_path" --query "SELECT 1"
rm -f "$init_db_client_config_path"