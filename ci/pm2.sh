#!/bin/bash

pids=(
    "client-dev             client:start"
    "backend-dev            backend:start"
    "redis                  redis"
    "logstash               logstash"
    "streaming-file         streaming:start:file"
    "streaming-api          streaming:start:api"
    "streaming-downloader   streaming:start:downloader"
    "streaming-encoder      streaming:start:encoder"
    "streaming-joiner       streaming:start:joiner"
)

function bunny {
    if which pm2 > /dev/null 2>&1; then
        for ((i=0; i<(${#pids[@]}); i++)); do
            echo "${pids[$i]}" | {
                while read name cmd; do
                    echo "pm2 start --name $name npm -- run $cmd";
                    # if ! pm2 ls | grep -q "$name"; then
                    #     pm2 start --name "$name" npm -- run $cmd;
                    # fi
                done;
            }
        done; return 1;
    else
        bunny
    fi
}
bunny
