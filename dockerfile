# Use a compatible base image with full build tools
FROM golang:1.22

# Enable CGO and set environment variables
ENV CGO_ENABLED=1 

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.* ./

# Download dependencies
RUN go mod download

# Install necessary build tools and libraries
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    sqlite3 && \
    apt-get clean

# Copy application source code
COPY . .

# Build the application
RUN go build -o main main.go

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]

