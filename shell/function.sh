#!/usr/bin/env bash

function add {
  echo "$(($1 + $2))"
}

echo "$(add 1 2)"