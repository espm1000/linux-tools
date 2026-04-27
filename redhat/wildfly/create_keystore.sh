#!/bin/bash

# Script to create quick and dirty keystore

if ! command -v keytool >/dev/null 2>&1; then
  echo "Keytool not found."
  exit 1
else
  echo "Keytool found..."
fi

echo "Creating keystore..."

keytool \
  -genkey \
  -alias server \
  -validity 365 \
  -keyalg RSA \
  -keystore keystore.jks \
  -dname "CN=localhost, OU=Servers, L=Home, ST=MN, C=US" \
  -storepass "password" >/dev/null 2>&1

  printf "Done.  Keystore created as \'keystore.jks\' with password of \'password\'\n"
