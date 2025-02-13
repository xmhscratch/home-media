#!/bin/bash

FFMPEG_INPUT_FILE=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")

cmd1=(
    "ffprobe"
    "-v error"
    "-select_streams a"
    "-show_streams"
    "-of default=noprint_wrappers=1:nokey=1"
    "-print_format json"
    "'$FFMPEG_INPUT_FILE'"
)

# echo ${cmd1[@]};
echo ${cmd1[@]} | sh | jq -cM '[.streams[] | {stream_index: .index, codec_name, lang_code: .tags.language, lang_title: .tags.title}]'
