# Use the official Go image as the base image
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go application with static linking for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY .env ./
EXPOSE 3000
CMD ["./server"]
