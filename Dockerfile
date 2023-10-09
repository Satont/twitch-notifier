FROM golang:alpine as builder

WORKDIR /app

RUN apk add --no-cache git curl wget upx make

COPY . .

RUN make build

FROM alpine:latest
RUN apk add --no-cache make curl && \
    curl -sSf https://atlasgo.sh | sh && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/build-out /app/
COPY --from=builder /app/docker-entrypoint.sh /app/
COPY --from=builder /app/Makefile /app/
COPY --from=builder /app/locales /app/locales
COPY --from=builder /app/ent/migrate /app/ent/migrate

RUN chmod +x /app/docker-entrypoint.sh

CMD ["sh", "/app/docker-entrypoint.sh"]
