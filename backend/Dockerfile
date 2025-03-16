FROM golang:1.22-bookworm

WORKDIR /app

# Copy the source code
COPY . .

# Create uploads directory with proper permissions
RUN mkdir -p /app/uploads && chmod 777 /app/uploads

# Install FFmpeg for video processing
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Expose port
EXPOSE 8080

# Run the application without building a binary
CMD ["go", "run", "./cmd/main.go"]
