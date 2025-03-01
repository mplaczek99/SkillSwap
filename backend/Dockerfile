# Build Stage
FROM golang:1.22-bookworm AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o SkillSwap cmd/main.go

# Final Stage
FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/SkillSwap .

# Create uploads directory
RUN mkdir -p /root/uploads && chmod 777 /root/uploads

# Install FFmpeg for video processing
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Configure environment variables
ENV CORS_ALLOW_ALL=true

EXPOSE 8080
CMD ["./SkillSwap"]
