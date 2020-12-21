FROM golang:1.13

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /data/gshark
WORKDIR $GOPATH/src/github.com/madneal/gshark
COPY . $GOPATH/src/github.com/madneal/gshark
RUN go build main.go

EXPOSE 8000
CMD ./main $OPTION