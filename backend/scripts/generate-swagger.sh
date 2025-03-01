#!/bin/bash
# File: ./backend/scripts/generate-swagger.sh

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Navigate to the backend directory
cd "$(dirname "$0")/.." || exit

# Remove old docs
rm -rf ./docs

# Generate new docs
echo "Generating Swagger documentation..."
swag init -g cmd/main.go -o docs

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html"
