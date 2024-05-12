#!/usr/bin/env bash
set -x

LOG_NAME=${1:-'test.log'}

./bin/syz-manager -config=syzkaller_docker.cfg 2>&1 | tee $LOG_NAME
