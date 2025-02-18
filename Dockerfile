# Build Stage
FROM golang:1.24.0-bookworm AS builder
WORKDIR /app

# Install git using apt-get (Debian's package manager)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Cache and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o SkillSwap cmd/main.go

# Final Stage
FROM debian:bookworm-slim
WORKDIR /root/

COPY --from=builder /app/SkillSwap .

EXPOSE 8080
CMD ["./SkillSwap"]
