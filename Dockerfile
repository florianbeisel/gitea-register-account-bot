# Start from a Go base image
FROM golang:latest as builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o mybot .

# Use a small base image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/mybot .

# Command to run the executable
CMD ["./mybot"]
