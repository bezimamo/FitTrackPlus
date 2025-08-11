@echo off
echo 🔧 Generating Swagger documentation...

REM Install swag if not installed
where swag >nul 2>nul
if %errorlevel% neq 0 (
    echo 📦 Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate docs
swag init -g cmd/server/main.go -o docs

echo ✅ Swagger documentation generated successfully!
echo 📚 Visit: http://localhost:8080/swagger/index.html
pause 