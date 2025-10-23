# Build stage
FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server .

# Runtime stage
FROM debian:bullseye-slim

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Download and install ONNX Runtime
RUN wget https://github.com/microsoft/onnxruntime/releases/download/v1.18.0/onnxruntime-linux-x64-1.18.0.tgz && \
    tar -xzf onnxruntime-linux-x64-1.18.0.tgz && \
    cp onnxruntime-linux-x64-1.18.0/lib/* /usr/local/lib/ && \
    ldconfig && \
    rm -rf onnxruntime-linux-x64-1.18.0*

# Copy binary from builder
COPY --from=builder /app/server .

# Create models directory
RUN mkdir -p models

# Expose port
EXPOSE 8090

# Run the application
CMD ["./server", "serve"]

#test@gmail.com
#password