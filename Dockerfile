FROM golang:1.20.0-alpine3.17 AS builder

ENV GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 

WORKDIR /simple_bank

COPY . .

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN go build -o main main.go 

FROM alpine:3.17 AS final

WORKDIR /simple_bank

COPY .env .
COPY --from=builder /simple_bank/main .
COPY --from=builder /simple_bank/db/migration ./db/migration

EXPOSE 8080

ENV BASE_DIR=/simple_bank

# ENTRYPOINT [ "/simple_bank/main"] 