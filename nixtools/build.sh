#!/bin/bash

set -e


function buildDebian() {
    if [[ ! -f "Dockerfile-debian" ]]; then
      echo "Dockerfile not found for debian."
      exit 1
    fi
    docker build -f Dockerfile-debian -t debian .
}

function buildRedhat() {
    if [[ ! -f "Dockerfile-redhat" ]]; then
      echo "Dockerfile not found for redhat."
      exit 1
    fi
    docker build -f Dockerfile-redhat -t redhat .
}


function main() {
    buildDebian
    buildRedhat
}

main
