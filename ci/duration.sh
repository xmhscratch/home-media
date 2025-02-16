#!/bin/bash

FFMPEG_INPUT_FILE=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")

# ../ci/duration.sh "$FFMPEG_INPUT_FILE";

cmd1=(
    "ffprobe"
    "-v error"
    "-show_entries format=duration"
    "-of default=noprint_wrappers=1:nokey=1"
    "'$FFMPEG_INPUT_FILE'"
)

# echo ${cmd1[@]};
echo ${cmd1[@]} | sh;
