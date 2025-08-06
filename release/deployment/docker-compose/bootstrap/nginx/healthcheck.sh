#!/bin/sh

set -e

curl -s http://localhost:80 | grep -E -q 'cozeloop|nginx'