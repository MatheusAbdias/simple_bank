FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

FROM scratch AS final
WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080
ENTRYPOINT [ "/app/main" ] 