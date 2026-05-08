#!/bin/bash

set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"


TOOL_FILENAME="tools.dev"
SECONDS=0

function cleandocker() {
  ~/bin/cleandocker
}

function buildTool() {
  if [[ ! -f "${TOOL_FILENAME}" ]]; then
    echo -e "${RED}binary not found, building...${DEFAULT}"
    go build -o "${TOOL_FILENAME}" .
  else
    echo -e "${RED}previous binary found, removing${DEFAULT}"
    rm -f "${TOOL_FILENAME}"
    go build -o "${TOOL_FILENAME}" .
  fi
}


function buildDebian() {
    if [[ ! -f "Dockerfile-debian" ]]; then
      echo -e "${RED}Dockerfile not found for Debian.${DEFAULT}"
      exit 1
    fi
    echo -e "${GREEN}Building Debian test image...${DEFAULT}"
    docker build -f Dockerfile-debian -t debian .
}

function buildRedhat() {
    if [[ ! -f "Dockerfile-redhat" ]]; then
      echo -e "${RED}Dockerfile not found for Redhat.${DEFAULT}"
      exit 1
    fi
    echo -e "${GREEN}Building RedHat test image...${DEFAULT}"
    docker build -f Dockerfile-redhat -t redhat .
}


function main() {
    if [[ $1 == "--clean" ]]; then
      cleandocker
      buildTool
      buildDebian
      buildRedhat
    else
      buildTool
      buildDebian
      buildRedhat
  fi
}

main "$@" 
echo -e "${GREEN}Elapsed time: $SECONDS seconds.${DEFAULT}"