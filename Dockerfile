# Start from a Go base image
FROM golang:alpine3.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY go.mod .
COPY go.sum .

# Download required modules
RUN go mod download 

# Copy the main application file
COPY main.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gitea-register-account-bot .

# Use a small base image
FROM alpine:edge

# Create and set the application directory
WORKDIR /app/

# Add a non-root user to run the application
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Copy the binary from the builder stage
COPY --from=builder /app/gitea-register-account-bot /app/

# Change file ownership to the nonroot user
RUN chown -R nonroot:nonroot /app

# Change to nonroot user
USER nonroot

# Command to run the executable
CMD ["./gitea-register-account-bot"]
