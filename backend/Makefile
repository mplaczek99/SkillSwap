.PHONY: run build test swagger clean

# Default target
all: swagger build

# Run the application
run:
	go run cmd/main.go

# Build the application
build:
	go build -o bin/skillswap cmd/main.go

# Run tests
test:
	go test -v ./...

# Generate Swagger documentation
swagger:
	bash scripts/generate-swagger.sh

# Clean generated files
clean:
	rm -rf bin
	rm -rf docs
