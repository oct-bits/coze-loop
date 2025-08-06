#!/bin/bash

set -e

client_conf() {
  cat <<EOF
$(eval "echo \"$(cat "/etc/cozeloop-mysql/client.cfg")\"")
EOF
}

exec mysqladmin --defaults-extra-file=<(client_conf) ping --silent