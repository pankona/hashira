#!/bin/bash -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

PROJECT_ID="hashira-auth"

IMPORT_FILE_ABSDIR="$(cd $(dirname $1) && pwd)"
IMPORT_FILE_BASENAME="$(basename $1)"
IMPORT_FILE="${IMPORT_FILE_ABSDIR}/${IMPORT_FILE_BASENAME}"

if [ ! -f "$IMPORT_FILE" ]; then
    echo "${IMPORT_FILE} is not file. Please specify a file exported from datastore"
    exit 1
fi

curl -X POST localhost:8081/v1/projects/${PROJECT_ID}:import \
-H 'Content-Type: application/json' \
-d '{"input_url":"'${IMPORT_FILE}'"}'
