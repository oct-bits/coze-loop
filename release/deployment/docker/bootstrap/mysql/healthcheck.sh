#!/bin/sh

set -e

export MYSQL_PWD="${COZE_LOOP_MYSQL_PASSWORD}"
exec mysqladmin \
      -h coze-loop-mysql \
      -u "${COZE_LOOP_MYSQL_USER}" \
      ping \
      --silent