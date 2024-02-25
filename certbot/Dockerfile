# Dockerfile
FROM ubuntu:22.04

# читаем переменные среды из .env
ARG DOMAIN_EMAIL
ARG DOMAIN_URL

# устанавливаем переменные среды в переменные
ENV DOMAIN_EMAIL=$DOMAIN_EMAIL
ENV DOMAIN_URL=$DOMAIN_URL

WORKDIR /certbot
COPY . /certbot
WORKDIR /certbot

RUN apt-get update
RUN apt-get -y install certbot

# запускаем скрипт генерации
CMD ["sh", "gen-cert.sh"]
