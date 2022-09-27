#!/usr/bin/env bash

# basic operators
A=3
B=$((100 * $A + 5))
echo $B # => 305

# string operations
## Length
STRING="this is a string"
echo "length STRING: ${#STRING}" # => length STRING: 16
## Extraction
POS=0
LEN=4
echo ${STRING:$POS:$LEN} # => this
echo ${STRING:10}        # => this
echo ${STRING:(-6)}      # => this
## Replacement
REPLACE_STRING="to be or not to be"
echo "Replace first:      ${REPLACE_STRING[*]/be/eat}"        # => to eat or not to be
echo "Replace all:        ${REPLACE_STRING[*]//be/eat}"       # => to eat or not to eat
echo "Delete:             ${REPLACE_STRING[*]// not/}"        # => to be or to be
echo "Replace beginning:  ${REPLACE_STRING[*]/%to be/to eat}" # => to be or not to eat
echo "Replace end:        ${REPLACE_STRING[*]/#to be/to eat}" # => to eat or not to be

# file test operators
if [ -e "function.sh" ]; then
  echo "function.sh is exist!"
fi
