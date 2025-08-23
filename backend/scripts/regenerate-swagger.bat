@echo off
echo 🔄 Regenerating Swagger documentation...

REM Install swag if not already installed
swag --version >nul 2>&1
if errorlevel 1 (
    echo 📦 Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Regenerate Swagger docs
echo 📝 Generating Swagger docs...
swag init -g cmd/server/main.go -o docs

if errorlevel 0 (
    echo ✅ Swagger documentation regenerated successfully!
    echo 📚 You can now access the updated Swagger UI at: http://localhost:8080/swagger/index.html
) else (
    echo ❌ Failed to regenerate Swagger documentation
    pause
    exit /b 1
)

