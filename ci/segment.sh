#!/bin/bash

FFMPEG_INPUT_FILE=$1
FFMPEG_START_TIME=$2
FFMPEG_DURATION=$3
FFMPEG_OUTPUT_FILE=$4

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
    "-i '$FFMPEG_INPUT_FILE'"
    "-map 0:s:0"
    "-scodec webvtt"
    "'/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/output.vtt'"
)

echo ${cmd1[@]};
echo ${cmd1[@]} | sh;

cmd2=(
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
    "-b:a 96k"
    "-ac 2"
    "'$FFMPEG_OUTPUT_FILE.mp4'"
);

echo ${cmd2[@]};
echo ${cmd2[@]} | sh;
