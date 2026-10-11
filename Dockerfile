# syntax=docker/dockerfile:1

# Builder Image
FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -v -o ./bin/goilerplate ./cmd/goilerplate

# Distribution Image
FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /app/bin/goilerplate /goilerplate

ENTRYPOINT ["/goilerplate"]
