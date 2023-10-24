FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR /app

COPY . .

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN apk update \
    && apk add --no-cache\
    curl=8.4.0-r0 \
    && rm -rf /var/cache/apk/* \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.17 AS final
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]