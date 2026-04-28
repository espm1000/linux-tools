#!/bin/bash
# shellcheck source=/dev/null

set -e

. /etc/os-release

LIBS=-lpthread
export LIBS

if [[ ! "${ID_LIKE}" =~ (rhel|centos|fedora)$ ]]; then
  echo "not rpm based"
  exit 1
fi

yum install -y \
    binutils \
    gcc \
    libstdc++ \
    libXext \
    libXrender \
    libXtst

# APR
APR="apr-1.7.5"
APR_UTIL="apr-util-1.6.3"
PCRE="pcre2-10.45"
curl -LO "https://dlcdn.apache.org//apr/${APR}.tar.gz"
curl -LO "https://dlcdn.apache.org//apr/${APR_UTIL}.tar.gz"
curl -LO "https://github.com/PCRE2Project/pcre2/releases/download/pcre2-10.45/${PCRE}.tar.gz"

tar zxf "${APR}.tar.gz"
tar zxf "${APR_UTIL}.tar.gz"
tar zxf "${PCRE}.tar.gz"

# Apache
APACHE="httpd-2.4.63"

curl -LO "https://dlcdn.apache.org/httpd/${APACHE}.tar.gz"

tar zxf ${APACHE}.tar.gz

mkdir ${APACHE}/srclib/apr
mkdir ${APACHE}/srclib/apr-util

cp -r ${APR}/* ${APACHE}/srclib/apr
cp -r ${APR_UTIL}/* ${APACHE}/srclib/apr-util

cd ${PCRE}/
./configure --prefix=/usr/local/pcre
make -j4
make -j4 install

cd ../${APACHE}/
./configure --enable-module=so --with-included-apr --with-pcre=/usr/local/pcre/bin/pcre2-config
