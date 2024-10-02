FROM golang:1.23 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY main.go ./main.go

RUN go build -v -o service-function-server

FROM ubuntu:24.04
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/service-function-server /app/service-function-server
CMD ["/app/service-function-server"]
