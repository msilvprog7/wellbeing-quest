FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/service

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/service .
COPY .env .
COPY cmd/service/*.sql ./cmd/service/

EXPOSE 8080
CMD ["./service"]