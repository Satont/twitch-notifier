FROM alpine:latest as builder

COPY --from=golang:alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add --no-cache git curl wget upx make

COPY libs libs
COPY go.mod go.sum /app/
RUN go mod download

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

RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/build-out"]
