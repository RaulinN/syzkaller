#!/usr/bin/env bash
set -x

# Prerequisites : currently in a syzkaller repo

TAGS="profiling"
GOTAGS="GOTAGS=\"${TAGS}\""

eval "make clean ${GOTAGS}"
eval "make generate ${GOTAGS}"
eval "make ${GOTAGS}"
