FROM golang:1.25.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-gateway ./cmd/main.go

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api-gateway /app/api-gateway
COPY --from=builder /app/internal/config/config.yaml /app/internal/config/config.yaml

EXPOSE 8080

ENTRYPOINT ["/app/api-gateway"]
