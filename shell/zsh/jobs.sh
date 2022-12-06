#!/usr/bin/env zsh

# protofmt
# shellcheck disable=SC2034
protofmt_content='
# protofmt buf format ./*.proto file
function protofmt() {
  for filename in $(find . -name "*.proto"); do
    buf format $filename -o $filename;
  done
}'

# gnum
# shellcheck disable=SC2034
gnum_content='
# gnum generate every N..M num on each line
function gnum() {
  echo {$1..$2} | tr " " "\n" | pbcopy
}'
