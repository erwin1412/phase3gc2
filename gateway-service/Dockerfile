# ========================
# Stage 1: Build the binary
# ========================
FROM golang:1.24-alpine AS build

# Install git (kadang perlu untuk go mod download)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (untuk cache dependency layer)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary (pastikan path-nya benar!)
RUN go build -o gateway ./cmd/main.go

# ===========================
# Stage 2: Run minimal image
# ===========================
FROM alpine:latest

# Install CA certificates for HTTPS request proxying
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from build stage
COPY --from=build /app/gateway .
COPY --from=build /app/docs ./docs

# ✅ Copy .env juga
COPY .env .

# Expose port Gateway
EXPOSE 8084

# Command to run
CMD ["./gateway"]
