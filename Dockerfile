FROM golang:1.21-alpine3.18 AS builder

WORKDIR /build

ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go mod download

RUN go build -o webhook github.com/benkenobi3/dick-and-dot/cmd/webhook
RUN go build -o longpool github.com/benkenobi3/dick-and-dot/cmd/longpool

FROM alpine:3.18

WORKDIR /bot

COPY --from=builder /build/webhook webhook
COPY --from=builder /build/longpool longpool
