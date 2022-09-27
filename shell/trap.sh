#!/usr/bin/env bash

#trap "echo Booh!" SIGINT SIGSTOP
trap "echo Booh!" SIGINT
echo "it's going to run until you hit Ctrl+Z"
echo "hit Ctrl+C to be blown away!"

while true
do
    sleep 60
done