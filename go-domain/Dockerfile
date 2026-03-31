# AWS Public ECR mirror of Docker Hub official images (more reliable pulls in CI)
FROM public.ecr.aws/docker/library/golang:1.25-alpine

# Install required packages (wget for Docker healthchecks in compose)
RUN apk add --no-cache git wget

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
