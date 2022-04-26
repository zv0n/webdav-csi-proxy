FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true
RUN apt update ; apt install -y git golang

WORKDIR /go/src/github.com/zv0n/webdav-proxy

ADD . /go/src/github.com/zv0n/webdav-proxy

RUN go mod download && go build -o /webdav-proxy
