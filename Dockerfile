FROM golang:1.24-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o leaflang ./cmd/leaflang

FROM alpine:latest

WORKDIR /

COPY --from=builder /leaflang /leaflang
COPY --from=builder /internal /internal
COPY --from=builder /pkg /pkg

EXPOSE 8080

CMD ["/leaflang"]
