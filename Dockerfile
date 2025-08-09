# Builder Image
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -v -o ./bin/ .

# Distribution Image
FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/bin/goilerplate /usr/bin/goilerplate

ENTRYPOINT ["/usr/bin/goilerplate"]
