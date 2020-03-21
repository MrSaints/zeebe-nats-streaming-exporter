#!/bin/bash

set -e

docker network create zeebe-dev-net || true

docker run --rm --name zeebe-dev-db \
    -p 5432:5432 \
    --env-file .dev.env \
    --network zeebe-dev-net \
    -v $(pwd)/_pgdata/:/var/lib/postgresql/data/:cached \
    postgres:11.1-alpine
