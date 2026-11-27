# syntax=docker/dockerfile:1

FROM oven/bun:1.4.2 AS frontend

WORKDIR /app

COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile

COPY frontend/ ./
RUN bun run build

FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
COPY --from=frontend /app/dist ./frontend/dist
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -v -o ./bin/goilerplate ./cmd/goilerplate

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /app/bin/goilerplate /goilerplate

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/goilerplate"]
CMD ["serve", "--port", "8080"]
