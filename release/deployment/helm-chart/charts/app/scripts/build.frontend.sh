#!/bin/sh

set -ex

OUTPUT_DIR="/cozeloop-bin/frontend/dist"

sh /cozeloop/frontend/apps/cozeloop/build-artifact.sh ${OUTPUT_DIR}
ls -lh "${OUTPUT_DIR}/"