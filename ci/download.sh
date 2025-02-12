#!/bin/bash

DOWNLOAD_URL=$1
OUTPUT_DIR=$2
BASE_URL=$3
ROOT_DIR=$4

# echo "$DOWNLOAD_URL"
# echo "$OUTPUT_DIR"
# echo "$BASE_URL"
# echo "$ROOT_DIR"

# echo "$ROOT_DIR/$OUTPUT_DIR"
# echo "$BASE_URL/$DOWNLOAD_URL"

# ../ci/download.sh "$DOWNLOAD_URL" "$OUTPUT_DIR" "$BASE_URL" "$ROOT_DIR";

cmd=(
  "wget2"
  "--no-check-certificate"
  "--quiet"
  "--method GET"
  "--timeout=30"
  "--header 'Content-Type: text/plain'"
  "--output-document '$ROOT_DIR/$OUTPUT_DIR'"
  "'$BASE_URL/$DOWNLOAD_URL'"
)

echo ${cmd[@]};
echo ${cmd[@]} | sh;
