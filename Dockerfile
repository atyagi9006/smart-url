# Use minimal Go image
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o url-shortener ./cmd

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./url-shortener"]
