FROM node:12.11.1-alpine

RUN apk add --no-cache bash

EXPOSE 3000

COPY . /app
WORKDIR /app

RUN npm i
RUN npm run build
RUN cp database.js.example database.js
RUN cp config.js.example config.js

CMD npm start