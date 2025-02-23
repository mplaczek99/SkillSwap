# Build Stage
FROM golang:1.24.0-bookworm AS builder
WORKDIR /app

# Install git (needed by go mod)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Cache and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire codebase and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o SkillSwap cmd/main.go

# Final Stage
FROM debian:bookworm-slim
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/SkillSwap .

# Expose the default port (for clarity)
EXPOSE 8080

# Run the compiled Go binary
CMD ["./SkillSwap"]

