# Use the official Go image as the base
FROM golang:latest AS builder

WORKDIR /app

# Copy the application code into the container
COPY . .

# Build the Go application
RUN go build -o pokedex-app

# Create a minimal runtime image
FROM debian:buster-slim
WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/pokedex-app .

# Command to run the application
CMD ["./pokedex-app"]