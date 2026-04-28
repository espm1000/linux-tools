FROM fedora:45

RUN dnf update -y
RUN dnf install -y \
    openssl

COPY ["create_ca.sh", "create_csr.sh", "openssl.cnf", "/scripts/"]