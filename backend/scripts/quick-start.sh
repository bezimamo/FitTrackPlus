#!/bin/bash

# FitTrack+ Quick Start Script
# This script helps you quickly set up and run the FitTrack+ backend

echo "🚀 FitTrack+ Backend Quick Start"
echo "================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Download from: https://golang.org/dl/"
    exit 1
fi

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install PostgreSQL first."
    echo "   Windows: https://www.postgresql.org/download/windows/"
    echo "   macOS: brew install postgresql"
    echo "   Linux: sudo apt-get install postgresql"
    exit 1
fi

echo "✅ Prerequisites check passed!"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please update it with your database credentials."
else
    echo "✅ .env file already exists."
fi

# Check if database exists
echo "🔍 Checking database connection..."
if psql -U postgres -d fittrackplus -c "SELECT 1;" &> /dev/null; then
    echo "✅ Database connection successful!"
else
    echo "⚠️  Database connection failed. Please run the database setup script:"
    echo "   psql -U postgres -f scripts/setup-database.sql"
    echo ""
    echo "   Or manually create the database:"
    echo "   CREATE DATABASE fittrackplus;"
    echo "   CREATE USER fittrackplus_user WITH PASSWORD 'fittrackplus_password_2024';"
    echo "   GRANT ALL PRIVILEGES ON DATABASE fittrackplus TO fittrackplus_user;"
fi

echo ""
echo "🎯 Ready to start the server!"
echo "   Run: go run cmd/server/main.go"
echo ""
echo "📚 Access points:"
echo "   - API Documentation: http://localhost:8080/api/docs"
echo "   - Swagger UI: http://localhost:8080/swagger/index.html"
echo "   - Health Check: http://localhost:8080/api/v1/health"
echo ""
echo "🧪 Test the API:"
echo "   curl -X GET http://localhost:8080/api/v1/health" 