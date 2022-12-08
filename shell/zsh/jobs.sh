#!/usr/bin/env zsh

# protofmt
# shellcheck disable=SC2034
protofmt_content='
# protofmt buf format ./*.proto file
function protofmt() {
  buf --version > /dev/null 2>&1
  if [ $? -ne 0 ]; then
    echo 1>&2 "protofmt: command not found: buf"
    echo 1>&2 "  (run: go install github.com/bufbuild/buf/cmd/buf@latest to install)";
    return;
  fi
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

# iterm
# shellcheck disable=SC2034
iterm_content='
# iterm launch new iterm2 tab at specail path
function iterm() {
  open -a iterm $1
}'
