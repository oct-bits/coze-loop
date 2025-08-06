#!/bin/sh

set -e

curl -f "http://localhost:${COZELOOP_MINIO_PORT_POD}/minio/health/live"