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

client_conf() {
  cat <<EOF
$(eval "echo \"$(cat "/etc/cozeloop-mysql/client.cfg")\"")
EOF
}

export MYSQL_ROOT_PASSWORD=${COZELOOP_MYSQL_PASSWORD}
export MYSQL_DATABASE=${COZELOOP_MYSQL_DATABASE}

print_banner "Mysql: Starting..."
print_banner_delay "Mysql: Successfully Started!" 12

echo "+ docker-entrypoint.sh mysqld"
docker-entrypoint.sh mysqld &

until mysqladmin --defaults-extra-file=<(client_conf) ping --silent; do
  sleep 2
done

i=1
# shellcheck disable=SC2010
for file in $(ls /etc/cozeloop-mysql/init-sql | grep '\.sql$'); do
  echo "+ init #$i: mysql -D $COZELOOP_MYSQL_DATABASE < $file"
  mysql --defaults-extra-file=<(client_conf) -D "$COZELOOP_MYSQL_DATABASE" < "/etc/cozeloop-mysql/init-sql/$file"
  i=$((i + 1))
done

wait