#!/usr/bin/env bash

function out() {
  echo "Standard output"
}

function err() {
  echo "Standard error" 1>&2
}

function isNull() {
  read in
  if [[ ${#in} -eq 0 ]]; then
    echo "True"
  else
    echo "False"
  fi
}

echo "out stdout: $(out 2>/dev/null)"
echo "err stdout: $(err 2>/dev/null)"
echo
echo "out isNull: $(out 2>/dev/null | isNull)"
echo "err isNull: $(err 2>/dev/null | isNull)"
echo "err redirect isNull: $(err 2>&1 | isNull)"
