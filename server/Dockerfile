FROM golang:1.18 as builder

WORKDIR /go/src/github.com/madneal/gshark/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="root@madneal.com"

WORKDIR /go/src/github.com/madneal/gshark/server

COPY --from=0 /go/src/github.com/madneal/gshark/server/server ./
COPY --from=0 /go/src/github.com/madneal/gshark/server/resource ./resource/
COPY --from=0 /go/src/github.com/madneal/gshark/server/config.docker.yaml ./

EXPOSE 8888
ENTRYPOINT ./server -c config.docker.yaml