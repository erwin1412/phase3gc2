# Gunakan base image minimal
FROM golang:1.24-alpine AS build

# Buat direktori kerja
WORKDIR /app

# Copy go.mod & go.sum dulu untuk caching deps
COPY go.mod go.sum ./

# Download deps
RUN go mod download

# Copy semua source code
COPY . .

# Build binary, misal jadi 'auth-service'
RUN go build -o auth-service ./cmd/main.go

# Expose HTTP & gRPC port
EXPOSE 8085
EXPOSE 50052

# Jalankan binary
CMD ["./auth-service"]
