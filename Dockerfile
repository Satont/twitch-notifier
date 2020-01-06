FROM node:12.11.1-alpine

ENV DB_NAME DB_NAME
ENV DB_USER DB_USER
ENV DB_PASSWORD DB_PASSWORD
ENV DB_HOST DB_HOST
ENV DB_PORT DB_PORT
ENV VKTOKEN VKTOKEN
ENV TWITCH_CLIENTID TWITCH_CLIENTID

RUN apk add --no-cache bash

EXPOSE 3000

COPY . /app
WORKDIR /app

RUN npm i
RUN npm run build
RUN cp database.js.example database.js
RUN cp config.js.example config.js

CMD npm start