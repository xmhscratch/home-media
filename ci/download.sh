#!/bin/bash

DOWNLOAD_URL=$(echo $1 | sed "s/^['\"]//; s/['\"]$//")
OUTPUT_DIR=$(echo $2 | sed "s/^['\"]//; s/['\"]$//")
BASE_URL=$(echo $3 | sed "s/^['\"]//; s/['\"]$//")
ROOT_DIR=$(echo $4 | sed "s/^['\"]//; s/['\"]$//")

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
