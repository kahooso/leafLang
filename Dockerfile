FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o leaflang ./cmd/leafLang

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/leaflang .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/pkg ./pkg

EXPOSE 8080

CMD ["./leaflang"]
