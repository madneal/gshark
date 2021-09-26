FROM golang:alpine as builder

WORKDIR /go/src/github.com/madneal/gshark/server
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=0
RUN go env
RUN go mod tidy
RUN go build -o server .

FROM alpine:latest
ENV GSHARK_CONFIG=/go/src/github.com/madneal/gshark/server/config.docker.yaml

WORKDIR /go/src/github.com/madneal/gshark/server

COPY --from=0 /go/src/github.com/madneal/gshark/server ./

EXPOSE 8888

ENTRYPOINT ./server web
