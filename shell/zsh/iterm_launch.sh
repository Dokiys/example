#!/usr/bin/env bash

# 导入job_builder.sh
. ./job_builder.sh --source-only

CONTENT='
# 整理proto
function protofmt4() {
  for filename in $(find . -name "*.proto"); do
    buf format $filename -o $filename;
  done
}'

job::build insert "$CONTENT"