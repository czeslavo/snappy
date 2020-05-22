FROM golang:alpine3.11 as builder

COPY . /build/
WORKDIR /build
RUN go build -o snappy ./cmd/snappy/main.go

FROM alpine:3.11

COPY --from=builder /build/snappy /app/snappy
VOLUME /app/snapshots

CMD ["/app/snappy", "-dir", "/app/snapshots"]
