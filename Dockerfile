# Start from a Go base image
FROM golang:alpine3.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .
RUN go mod download 

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mybot .

# Use a small base image
FROM alpine:edge

WORKDIR /app/

# Copy the binary from the builder stage
COPY --from=builder /app/mybot /app/

# Command to run the executable
CMD ["./mybot"]
