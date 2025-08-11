# FitTrack+ API Testing Commands

## 🚀 Quick Test Commands

### 1. Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### 2. Register a New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "member"
  }'
```

### 3. Login User
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 4. Get User Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### 5. Update User Profile (Protected)
```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "first_name": "John Updated",
    "last_name": "Doe Updated",
    "phone": "+0987654321"
  }'
```

## 🔧 PowerShell Commands (Windows)

### 1. Health Check
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health" -Method GET
```

### 2. Register User
```powershell
$body = @{
    email = "john@example.com"
    password = "password123"
    first_name = "John"
    last_name = "Doe"
    phone = "+1234567890"
    role = "member"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $body -ContentType "application/json"
```

### 3. Login User
```powershell
$body = @{
    email = "john@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -Body $body -ContentType "application/json"
$token = $response.token
```

### 4. Get Profile
```powershell
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/profile" -Method GET -Headers $headers
```

## 🧪 Automated Testing (Better Approach)

Instead of `test_auth.go`, you should use:

1. **Unit Tests** - Test individual functions
2. **Integration Tests** - Test API endpoints
3. **Swagger UI** - Interactive testing
4. **Postman/Insomnia** - API client tools

## 🎯 Why Not Use `test_auth.go`?

- **Not reusable** - Hardcoded test data
- **No proper assertions** - Just prints results
- **Manual execution** - Need to run manually
- **Not integrated** - Doesn't use Go's testing framework
- **Maintenance burden** - Need to update when API changes

## ✅ Recommended Testing Approach

1. **Development**: Use Swagger UI for quick tests
2. **Manual Testing**: Use curl/PowerShell commands
3. **Automated Testing**: Write proper Go tests
4. **API Client**: Use Postman or similar tools 