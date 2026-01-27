FROM golang:1.25-alpine

# Install required packages
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /app/bin/api ./cmd/api

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/bin/api"]
