# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/adserver

RUN go build -o adserver ./cmd/main.go

# -----
FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/adserver .

ENTRYPOINT ["/app/adserver"]
