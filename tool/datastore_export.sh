#!/bin/bash -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

PROJECT_ID="hashira-auth"
EXPORT_DIR="$SCRIPT_DIR"/datastore-export-`date "+%Y%m%d-%H%M%S"`

curl -X POST localhost:8081/v1/projects/${PROJECT_ID}:export \
    -H 'Content-Type: application/json' \
    -d '{"output_url_prefix":"'${EXPORT_DIR}'"}'
