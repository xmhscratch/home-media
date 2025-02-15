#!/bin/bash

FFMPEG_INPUT_FILE=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")
STREAM_INDEX=$2
LANG_CODE=$3
FFMPEG_OUTPUT_FILE=$(echo $4 | sed "s/^['\"]//; s/['\"]$//")

# ../ci/extract-vid.sh "$FFMPEG_INPUT_FILE" "$FFMPEG_OUTPUT_FILE";

cmd1=(
    "ffmpeg -y"
    "-i '$FFMPEG_INPUT_FILE'"
    "-map 0:v"
    "-vcodec copy"
    "-sn -dn -an"
    "-map_metadata -1"
    "'$FFMPEG_OUTPUT_FILE.mp4'"
)

echo ${cmd1[@]};
echo ${cmd1[@]} | sh;
