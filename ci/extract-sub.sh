#!/bin/bash

FFMPEG_INPUT_FILE=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")
STREAM_INDEX=$2
LANG_CODE=$3
FFMPEG_OUTPUT_FILE=$(echo $4 | sed "s/^['\"]//; s/['\"]$//")

# echo $FFMPEG_INPUT_FILE
# echo $STREAM_INDEX
# echo $LANG_CODE
# echo $FFMPEG_OUTPUT_FILE

# ../ci/extract-sub.sh "$FFMPEG_INPUT_FILE" "$STREAM_INDEX" "$LANG_CODE" "$FFMPEG_OUTPUT_FILE";

cmd1=(
    "ffmpeg -y"
    "-i '$FFMPEG_INPUT_FILE'"
    "-map 0:$STREAM_INDEX"
    "-scodec webvtt"
    "'$FFMPEG_OUTPUT_FILE.$LANG_CODE.vtt'"
)

echo ${cmd1[@]};
echo ${cmd1[@]} | sh;
