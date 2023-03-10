FROM registry-mirror.pingcap.net/library/debian:buster

ARG arch amd64

WORKDIR /apps/
RUN apt update \
    && apt install -y ca-certificates\
    && rm -rf /var/lib/apt/lists/*
ADD ./bin/linux/${arch} /apps/
ADD ./website/build/ /apps/website/build/
