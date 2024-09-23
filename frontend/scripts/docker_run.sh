#!/bin/bash

set -a
source ./.env
set +a

docker run --rm --name blog-frontend \
  --network test \
  -p 3000:3000 \
  "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG"
