# Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

# COPY go.mod go.sum ./
COPY go.mod ./
# skip go.sum because it doesn't exist
# RUN go mod download
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]