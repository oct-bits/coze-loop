#!/bin/bash

set -e

echo > "/dev/tcp/127.0.0.1/${COZELOOP_RMQNAMESRV_PORT_POD}" 2>/dev/null