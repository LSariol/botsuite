FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o botsuite ./cmd/botsuite

FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/botsuite

COPY --from=builder /app/botsuite /botsuite

EXPOSE 2400
CMD ["/botsuite"]