# Build Stage
FROM golang:1.24.0-bookworm AS builder
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
EXPOSE 8080
CMD ["./SkillSwap"]

