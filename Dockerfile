FROM node:15.13.0-alpine3.10

RUN apk add --no-cache bash

EXPOSE 3000

COPY . /app
WORKDIR /app

RUN npm i
RUN npm run build

CMD npm start