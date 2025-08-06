#!/bin/sh

set -e

wget -qO- http://localhost:8888/ping | grep -q pong