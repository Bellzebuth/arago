# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o adserver ./cmd/main.go

# -----
FROM debian:stable-slim

RUN apt-get update && apt-get install -y libc6

WORKDIR /app
COPY --from=builder /app/adserver .

ENTRYPOINT ["/app/adserver"]
