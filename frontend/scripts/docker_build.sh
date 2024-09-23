#!/bin/bash

set -a
source ./.env
set +a

docker build --network host -t "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG" \
  --build-arg BACKEND_BASE_URL=$BACKEND_BASE_URL .
