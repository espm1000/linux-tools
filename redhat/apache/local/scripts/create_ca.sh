#!/bin/bash

# This script will create a simple CA to be used for signing CSRs in a non-prod world.
set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"

# Vars
CERTS_FOLDER="../certs"

if ! command -v openssl >/dev/null 2>&1; then
  echo "OpenSSL not found"
else
  echo -e "${GREEN}OpenSSL installed...${DEFAULT}"
fi

if [[ ! -f "../conf/openssl.cnf" ]]; then
  echo -e "${RED}Missing openssl configuration file.${DEFAULT}"
  exit 1
fi

read -rp "Days: " DAYS
read -rp "CA Name: " CANAME

# No passphrase
openssl req -x509 \
  -noenc \
  -config ../conf/openssl.cnf \
  -days "${DAYS}" \
  -newkey rsa:4096 \
  -keyout "${CERTS_FOLDER}/${CANAME}-key.pem" \
  -out "${CERTS_FOLDER}/${CANAME}-cert.pem"
