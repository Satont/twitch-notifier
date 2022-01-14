FROM node:16-alpine

RUN apk add --no-cache bash


EXPOSE 3000
EXPOSE 9229

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN npm i -g pnpm && pnpm install

COPY . /app
RUN pnpm run build

COPY docker.sh /
RUN chmod +x /docker.sh
ENTRYPOINT ["/docker.sh"]