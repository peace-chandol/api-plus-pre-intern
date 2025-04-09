# Use the official Go image as the base for building the application
FROM golang:latest AS builder

WORKDIR /app

# Copy the application code into the container
COPY . .

# Build the Go application
RUN go build -o pokedex-app

# Create a minimal runtime image with a newer version of glibc
FROM debian:bookworm-slim

# Install necessary dependencies (glibc)
RUN apt-get update && apt-get install -y libc6

WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/pokedex-app .

# Command to run the application
CMD ["./pokedex-app"]
