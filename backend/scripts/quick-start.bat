@echo off
REM FitTrack+ Quick Start Script for Windows
REM This script helps you quickly set up and run the FitTrack+ backend

echo 🚀 FitTrack+ Backend Quick Start
echo =================================

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go first.
    echo    Download from: https://golang.org/dl/
    pause
    exit /b 1
)

REM Check if PostgreSQL is installed
psql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ PostgreSQL is not installed. Please install PostgreSQL first.
    echo    Download from: https://www.postgresql.org/download/windows/
    pause
    exit /b 1
)

echo ✅ Prerequisites check passed!

REM Install dependencies
echo 📦 Installing Go dependencies...
go mod tidy

REM Check if .env file exists
if not exist .env (
    echo 📝 Creating .env file from template...
    copy env.example .env
    echo ✅ .env file created. Please update it with your database credentials.
) else (
    echo ✅ .env file already exists.
)

echo.
echo 🎯 Ready to start the server!
echo    Run: go run cmd/server/main.go
echo.
echo 📚 Access points:
echo    - API Documentation: http://localhost:8080/api/docs
echo    - Swagger UI: http://localhost:8080/swagger/index.html
echo    - Health Check: http://localhost:8080/api/v1/health
echo.
echo 🧪 Test the API:
echo    curl -X GET http://localhost:8080/api/v1/health
echo.
pause 