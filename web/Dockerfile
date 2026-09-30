FROM --platform=$BUILDPLATFORM node:22-alpine@sha256:0a7108bf6c7bf5de370ffb1a3ed6be93d405b43ff159f681a8d18c0e2bc2e402 AS builder

WORKDIR /gshark_web/
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps
COPY . .
RUN npm run build

FROM nginx:stable-alpine@sha256:0985e772fb9f729e6fa0980da05fca5d9c468e870eed43071545afa9d2e27d94
LABEL MAINTAINER="root@madneal.com"

COPY my.conf /etc/nginx/conf.d/my.conf
COPY --from=builder /gshark_web/dist /usr/share/nginx/html
EXPOSE 8080
