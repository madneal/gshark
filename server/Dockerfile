# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25.14-alpine3.24 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -mod=readonly -trimpath -o /out/gshark .

FROM alpine:3.24.2
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /go/src/github.com/madneal/gshark/server
ENV GSHARK_CONFIG=config.docker.yaml
COPY --from=builder /out/gshark ./gshark
COPY resource/rbac_model.conf ./resource/rbac_model.conf
COPY config.docker.yaml ./config.docker.yaml
EXPOSE 8888
ENTRYPOINT ["./gshark"]
CMD ["serve"]
