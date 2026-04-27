#!/bin/bash

# This script will create a simple CA to be used for signing CSRs in a non-prod world.
set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"

CONF_FOLDER="."
CONF_FILE="openssl.cnf"

if ! command -v openssl >/dev/null 2>&1; then
  echo "OpenSSL not found"
else
  echo -e "${GREEN}OpenSSL installed...${DEFAULT}"
fi

if [[ ! -f "openssl.cnf" ]]; then
  echo -e "${RED}Missing openssl configuration file.${DEFAULT}"
  exit 1
fi

read -rp "Days: " DAYS
read -rp "CA Name: " CANAME

# To remove the passphase:
# OpenSSL 3.2+ use '-noenc'
# OpenSSL 1.0.2 use '-nodes'
openssl req -x509 \
  -config "${CONF_FILE}" \
  -days "${DAYS}" \
  -newkey rsa:4096 \
  -keyout "${CONF_FOLDER}/${CANAME}-key.pem" \
  -out "${CONF_FOLDER}/${CANAME}-cert.pem"


: << EOT
openssl x509 \
-req \
-in my_cert_req.pem \
-days 365 \
-CA ca_cert.pem \
-CAkey ca_private_key.pem \
-CAcreateserial \
-out my_signed_cert.pem
EOT
