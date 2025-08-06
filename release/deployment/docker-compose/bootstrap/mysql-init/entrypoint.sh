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

print_banner "Starting..."

export MYSQL_PWD="${COZE_LOOP_MYSQL_PASSWORD}"

i=1
# shellcheck disable=SC2010
for file in $(ls /coze-loop-mysql-init/bootstrap/init-sql | grep '\.sql$'); do
  echo "+ Init #$i: $file"
  mysql \
    -h coze-loop-mysql \
    -u "${COZE_LOOP_MYSQL_USER}" \
    -D "${COZE_LOOP_MYSQL_DATABASE}" \
    < "/coze-loop-mysql-init/bootstrap/init-sql/${file}"
  i=$((i + 1))
done

print_banner "Completed!"