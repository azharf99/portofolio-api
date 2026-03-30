# Stage 1: Build aplikasinya
FROM golang:1.26-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum, lalu download dependency
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi
COPY . .

# Build aplikasi Golang (CGO_ENABLED=0 memastikan binary mandiri tanpa library C)
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-api main.go

# Stage 2: Container final yang sangat ringan
FROM alpine:latest

WORKDIR /app

# Salin binary dari stage builder ke stage final ini
COPY --from=builder /app/portfolio-api .

# Buka port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./portfolio-api"]