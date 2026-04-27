#!/bin/bash

# This script will create a key and CSR using OpenSSL

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

KEY_SIZE=2048

function create_csr() {
  openssl req \
    -nodes \
    -newkey rsa:"${KEY_SIZE}" \
    -keyout "${CERTS_FOLDER}/private_key.key" \
    -out "${CERTS_FOLDER}/sign_this.csr"
}

function signit() {
  openssl x509 \
    -req \
    -in "${CERTS_FOLDER}/sign_this.csr" \
    -days 20 \
    -CA "${CERTS_FOLDER}/local-cert.pem" \
    -CAkey "${CERTS_FOLDER}/local-key.pem" \
    -CAcreateserial \
    -out "${CERTS_FOLDER}/signed.pem"
}

case $1 in
--csr)
  create_csr
  ;;
--sign)
  signit
  ;;
esac
