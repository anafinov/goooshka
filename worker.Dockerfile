FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker worker/cmd/worker/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /worker .
EXPOSE 8081
CMD ["./worker"]
