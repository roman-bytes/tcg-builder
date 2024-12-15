# Use the official Go image as the base image
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o server .

# Final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Copy the .env file
COPY .env ./

# Expose the port the app runs on
EXPOSE 3000

# Command to run the server
CMD ["go", "run", "main.go"]