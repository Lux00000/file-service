FROM golang:1.23.8 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o file-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/file-service .
COPY --from=builder /app/swagger ./swagger
COPY --from=builder /app/api/swagger/api/proto/*.swagger.json ./swagger/api/proto/

EXPOSE 50051 8080
CMD ["./file-service"]