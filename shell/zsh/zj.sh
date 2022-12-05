#!/usr/bin/env bash

# 导入zsh_job_builder.sh
. ./job_builder.sh --source-only

main() {
  job::build "$@"
}

main "$@"