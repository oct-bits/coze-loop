#!/bin/sh

set -e

curl -s "http://localhost:${COZELOOP_INGRESS_PORT_POD}" | grep -E -q 'cozeloop|nginx'