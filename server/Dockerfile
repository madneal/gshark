# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25-alpine@sha256:1ae0735f00daffa3aaf1363a5184c0d2dc55c78e3db4ec70241cdac97bf84b59 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -mod=readonly -trimpath -o /out/gshark .

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /go/src/github.com/madneal/gshark/server
ENV GSHARK_CONFIG=config.docker.yaml
COPY --from=builder /out/gshark ./gshark
COPY resource ./resource
COPY config.docker.yaml ./config.docker.yaml
EXPOSE 8888
ENTRYPOINT ["./gshark"]
CMD ["serve"]
