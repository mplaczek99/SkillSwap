# Build Stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
# Cache and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o SkillSwap cmd/main.go

# Final Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/SkillSwap .
EXPOSE 8080
CMD ["./SkillSwap"]

