FROM golang:1.21.5 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o distributed_kv_storage ./cmd/distributed_kv_storage

FROM golang:1.21.5-alpine3.18
WORKDIR /
RUN apk add git
COPY --from=builder /build/distributed_kv_storage /usr/bin/distributed_kv_storage

EXPOSE 8080
CMD ["distributed_kv_storage"]

