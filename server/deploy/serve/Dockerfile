FROM golang:1.20 as builder

ENV GSHARK_CONFIG=config.docker.yaml
WORKDIR /go/src/github.com/madneal/gshark/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o gshark .

FROM alpine:latest

LABEL MAINTAINER="root@madneal.com"

WORKDIR /go/src/github.com/madneal/gshark/server

COPY --from=0 /go/src/github.com/madneal/gshark/server/gshark ./
COPY --from=0 /go/src/github.com/madneal/gshark/server/resource ./resource/
COPY --from=0 /go/src/github.com/madneal/gshark/server/config.docker.yaml ./

EXPOSE 8888
ENTRYPOINT ./gshark serve