#!/bin/bash

set -e


function run() {
  echo "starting server..."
  docker run \
    -p "8082:8082" \
    -p "8080:8080" \
    -p "9443:9443" \
    --name apache_local \
    -d \
    apache:local
}

function build() {
  set +e
  echo "building..."
  docker stop apache_local
  docker container prune -f
  docker build -f ../Dockerfile ../ -t apache:local
}

case $1 in
--run)
  run
  ;;
--build)
  build
  ;;
*)
  run
  ;;
esac
