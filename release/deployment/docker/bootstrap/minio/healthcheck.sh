#!/bin/sh

set -e

curl \
  -f "http://coze-loop-minio:9000/minio/health/live" \
  > /dev/null 2>&1