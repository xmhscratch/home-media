#!/bin/bash

FFMPEG_INPUT_FILE=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")
FFMPEG_START_TIME=$(echo $2 | sed "s/^['\"]//; s/['\"]$//")
FFMPEG_DURATION=$(echo $3 | sed "s/^['\"]//; s/['\"]$//")
FFMPEG_OUTPUT_FILE=$(echo $4 | sed "s/^['\"]//; s/['\"]$//")

# echo $FFMPEG_INPUT_FILE
# echo $FFMPEG_START_TIME
# echo $FFMPEG_DURATION
# echo $FFMPEG_OUTPUT_FILE

# echo "../ci/segment.sh" "$FFMPEG_INPUT_FILE" "$FFMPEG_START_TIME" "$FFMPEG_DURATION" "$FFMPEG_OUTPUT_FILE" | sh

NVIDIA_SUPPORT=0
if ! command -v nvidia-smi &> /dev/null; then
    NVIDIA_SUPPORT=0
else
    if [[ -z "$(nvidia-smi -L)" ]]; then
        NVIDIA_SUPPORT=0
    else
        NVIDIA_SUPPORT=1
    fi;
fi;

cmd1=(
    "ffmpeg -y"
    "-ss '$FFMPEG_START_TIME'"
    $([[ -n "$NVIDIA_SUPPORT" && "$NVIDIA_SUPPORT" == "1" ]] && echo "-hwaccel cuda")
    "-i '$FFMPEG_INPUT_FILE'"
    "-t '$FFMPEG_DURATION'"
    "-threads $(getconf _NPROCESSORS_ONLN)"
    $([[ -n "$NVIDIA_SUPPORT" && "$NVIDIA_SUPPORT" == "1" ]] && echo "-c:v hevc_nvenc")
    $([[ -n "$NVIDIA_SUPPORT" && "$NVIDIA_SUPPORT" == "1" ]] && echo "-cq 22")
    $([[ -z "$NVIDIA_SUPPORT" || "$NVIDIA_SUPPORT" == "0" ]] && echo "-c:v libvpx-vp9")
    $([[ -z "$NVIDIA_SUPPORT" || "$NVIDIA_SUPPORT" == "0" ]] && echo "-preset slow")
    "-b:v 9000k"
    "-pix_fmt yuv420p"
    "-c:a libopus"
    "-b:a 128k"
    "-ac 2"
    "'$FFMPEG_OUTPUT_FILE.mp4'"
);

echo ${cmd1[@]};
echo ${cmd1[@]} | sh;
