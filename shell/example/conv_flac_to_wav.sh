#!/bin/zsh

oldext="flac"
newext="wav"

for file in *.${oldext}
do
  name=$(ls "$file" | cut -d. -f1)
  ffmpeg -i $file ${name}.$newext
done
