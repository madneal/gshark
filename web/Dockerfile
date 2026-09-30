FROM --platform=$BUILDPLATFORM node:22.23.3-alpine AS builder

WORKDIR /gshark_web/
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps
COPY . .
RUN npm run build

FROM nginx:1.30.5-alpine
LABEL MAINTAINER="root@madneal.com"

COPY my.conf /etc/nginx/conf.d/my.conf
COPY --from=builder /gshark_web/dist /usr/share/nginx/html
EXPOSE 8080
