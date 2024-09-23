#/bin/bash

set -a
source ./.env
set +a

docker run -d --rm --name blog-server \
  --network test \
  --mount type=bind,source="./config.json",target="/app/config.json" \
  --mount type=bind,source="./blog.db",target="/app/blog.db" \
  -p 8080:8080 \
  "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG" ./server --migrate
