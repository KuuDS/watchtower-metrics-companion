FROM golang:1.24-alpine AS builder

WORKDIR /app

# RUN apk add --no-cache git ca-certificates tzdata

# COPY go.mod go.sum ./

# RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o watchtower-metrics .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/watchtower-metrics .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./watchtower-metrics"]
