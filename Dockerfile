FROM node:15.13.0-alpine3.10

RUN apk add --no-cache bash

EXPOSE 3000
EXPOSE 9229

COPY . /app
WORKDIR /app

RUN npm i
RUN npm run build

COPY docker.sh /
RUN chmod +x /docker.sh
ENTRYPOINT ["/docker.sh"]