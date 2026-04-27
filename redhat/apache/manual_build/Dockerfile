FROM ubuntu:24.04

RUN apt update && apt upgrade -y

RUN apt install -y \
      libapr1 \
      libaprutil1 \
      libaprutil1-dev \
      libpcre2-dev \
      curl \
      build-essential

RUN mkdir /app
WORKDIR /app
RUN curl -LO https://dlcdn.apache.org/httpd/httpd-2.4.63.tar.gz
RUN tar zxf httpd-2.4.63.tar.gz
WORKDIR /app/httpd-2.4.63
RUN ./configure --prefix=/usr/local
