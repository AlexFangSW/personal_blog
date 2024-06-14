#/bin/bash

set -a
source ./.env
set +a

docker run --rm --name blog-server \
  --mount type=bind,source="./config.json",target="/app/config.json" \
  -p 8080:8080 \
  "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG" ./server --migrate
