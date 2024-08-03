FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN go build -o /sso-gc

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /sso-gc .

EXPOSE 8080

CMD ["./sso-gc"]