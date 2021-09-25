FROM node:12.16.1

WORKDIR /gshark_web/
COPY . .
RUN rm -rf package-lock.json
RUN npm install
RUN npm run build

FROM nginx:alpine

COPY .docker-compose/nginx/my.conf /etc/nginx/conf.d/my.conf
COPY --from=0 /gshark_web/dist /usr/share/nginx/html
RUN cat /etc/nginx/nginx.conf
RUN cat /etc/nginx/conf.d/my.conf
RUN ls -al /usr/share/nginx/html