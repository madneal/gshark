FROM node:16

WORKDIR /gshark_web/
COPY . .

RUN npm install --legacy-peer-deps && npm run build

FROM nginx:alpine
LABEL MAINTAINER="root@madneal.com"

COPY my.conf /etc/nginx/conf.d/my.conf
COPY --from=0 /gshark_web/dist /usr/share/nginx/html
RUN cat /etc/nginx/nginx.conf
RUN cat /etc/nginx/conf.d/my.conf
RUN ls -al /usr/share/nginx/html