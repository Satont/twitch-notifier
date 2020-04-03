FROM node:12.11.1-alpine

RUN apk add --no-cache bash

EXPOSE 3000

COPY . /app
WORKDIR /app

RUN npm i
RUN npm run build

CMD npm start