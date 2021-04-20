FROM golang:alpine

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/gshark
COPY . .
RUN go env && go build -o server .

FROM alpine:latest

WORKDIR /go/src/gshark
COPY --from=0 /go/src/github.com/madneal/gshark/server ./
COPY --from=0 /go/src/github.com/madneal/gshark/config.yaml ./
COPY --from=0 /go/src/github.com/madneal/gshark/resource ./resource

EXPOSE 8888

ENTRYPOINT ./server
