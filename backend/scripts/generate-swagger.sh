#!/bin/bash

# Generate Swagger documentation
echo "🔧 Generating Swagger documentation..."

# Install swag if not installed
if ! command -v swag &> /dev/null; then
    echo "📦 Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate docs
swag init -g cmd/server/main.go -o docs

echo "✅ Swagger documentation generated successfully!"
echo "📚 Visit: http://localhost:8080/swagger/index.html" 