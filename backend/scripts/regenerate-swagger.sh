#!/bin/bash

echo "🔄 Regenerating Swagger documentation..."

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "📦 Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Regenerate Swagger docs
echo "📝 Generating Swagger docs..."
swag init -g cmd/server/main.go -o docs

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation regenerated successfully!"
    echo "📚 You can now access the updated Swagger UI at: http://localhost:8080/swagger/index.html"
else
    echo "❌ Failed to regenerate Swagger documentation"
    exit 1
fi

