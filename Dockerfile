FROM golang:1.23.0-alpine3.20 AS builder

LABEL stage=builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN apk update --no-cache

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /build/webhook ./cmd/webhook
RUN go build -ldflags="-s -w" -o /build/longpool ./cmd/longpool

FROM alpine:3.20

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/webhook /app/webhook
COPY --from=builder /build/longpool /app/longpool

CMD ["./longpool"]
