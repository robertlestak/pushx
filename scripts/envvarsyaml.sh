#!/bin/bash

VARS=(`egrep -oh --exclude Makefile \
    --exclude-dir bin \
    --exclude-dir scripts \
    -R 'os.Getenv\(.*?\)' . | \
    tr -d ' ' | \
    sort | \
    uniq | \
    sed -e 's,os.Getenv(,,g' -e 's,),,g' \
    -e 's,",,g' \
    -e 's,prefix+,PUSHX_,g'`)
for VAR in ${VARS[@]}; do
    echo "- \`$VAR\`"
done
