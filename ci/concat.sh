#!/bin/bash

INPUT=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")
OUTPUT=$(echo $2 | sed "s/^['\"]//; s/['\"]$//")

# echo $INPUT
# echo $OUTPUT

# ../ci/concat.sh "$INPUT" "$OUTPUT";

cmd1=(
    "ffmpeg"
    "-f concat"
    "-safe 0"
    "-i '$INPUT'"
    "-c copy"
    "'$OUTPUT'"
)

echo ${cmd1[@]};
echo ${cmd1[@]} | sh;

# cmd2=(

# )

# echo ${cmd2[@]};
# echo ${cmd2[@]} | sh;
