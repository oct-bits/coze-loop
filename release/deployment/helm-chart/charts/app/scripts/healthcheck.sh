#!/bin/sh

curl -s "http://localhost:${COZELOOP_APP_HERTZ_PORT_POD}/ping" | grep -q pong