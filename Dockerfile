# syntax=docker/dockerfile:1

# --- build stage ---
FROM golang:1.26-alpine AS builder

ARG SWAG_VERSION=v1.16.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs so the image builds even if docs/ isn't committed.
RUN go run github.com/swaggo/swag/cmd/swag@${SWAG_VERSION} init \
    -g cmd/server/main.go --parseInternal --parseDependency -o docs

RUN CGO_ENABLED=0 go build -o /app/bin/api ./cmd/server

# --- runtime stage ---
FROM alpine:3.22

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/bin/api /app/api

USER app

EXPOSE 8080

ENTRYPOINT ["/app/api"]
