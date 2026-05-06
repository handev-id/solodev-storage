# syntax=docker/dockerfile:1.7

FROM golang:1.26.2-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/storage-api .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /out/storage-api /app/storage-api

ENV APP_PORT=3000
EXPOSE 3000

USER nonroot:nonroot
ENTRYPOINT ["/app/storage-api"]
