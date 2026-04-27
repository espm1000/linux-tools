FROM ubuntu:latest

RUN apt update && apt upgrade -y
RUN apt install -y \
    vim \
    apache2
COPY ./conf/ports.conf /etc/apache2/ports.conf
COPY ./conf/000-default.conf /etc/apache2/sites-enabled/hello.conf
COPY ./www/hello/ /var/www/hello/
COPY ./certs/private_key.key /etc/apache2/certs/private_key.key
COPY ./certs/signed.pem /etc/apache2/certs/signed.pem

# Start apache

CMD ["apachectl", "-k", "start", "-X"]
#CMD ["sleep", "900"]