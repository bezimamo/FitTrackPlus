# FitTrack+ Backend

A Go-based backend API for the FitTrack+ fitness platform.

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd FitTrackPlus/backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your database credentials
   ```

4. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE fittrackplus;
   CREATE USER fittrackplus_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE fittrackplus TO fittrackplus_user;
   ```

5. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

## 🧪 Testing the Setup

### Test without database
```bash
go run test_server.go
```
Visit: http://localhost:8080/test

### Test with database
```bash
go run cmd/server/main.go
```
Visit: http://localhost:8080/api/docs

### Swagger Documentation
```bash
# Generate Swagger docs (Windows)
scripts/generate-swagger.bat

# Generate Swagger docs (Linux/Mac)
chmod +x scripts/generate-swagger.sh
./scripts/generate-swagger.sh
```
Visit: http://localhost:8080/swagger/index.html

## 📁 Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── common/
│   │   ├── config/
│   │   │   └── config.go    # Configuration management
│   │   ├── database/
│   │   │   └── database.go  # Database connection
│   │   └── models/
│   │       └── user.go      # Database models
│   ├── auth/                # Authentication service (coming soon)
│   ├── user/                # User management (coming soon)
│   ├── plan/                # Plan generation (coming soon)
│   ├── booking/             # Booking system (coming soon)
│   ├── progress/            # Progress tracking (coming soon)
│   ├── payment/             # Payment processing (coming soon)
│   └── content/             # Content management (coming soon)
├── migrations/              # Database migrations
├── pkg/                     # Reusable packages
├── scripts/                 # Build and generation scripts
├── docs/                    # Swagger documentation
├── go.mod                   # Go module file
├── go.sum                   # Dependency checksums
├── env.example              # Environment variables template
└── README.md               # This file
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `fittrackplus` |
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

## 📚 Learning Go Concepts

### 1. **Packages and Imports**
```go
package main  // Every executable Go program must have a main package

import (
    "fmt"     // Standard library for formatting
    "log"     // Standard library for logging
    "net/http" // Standard library for HTTP
)
```

### 2. **Structs and Methods**
```go
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Email string `json:"email" gorm:"uniqueIndex"`
}

// Method on User struct
func (u *User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}
```

### 3. **Pointers**
```go
// & gets the memory address
// * dereferences a pointer
user := &User{Email: "test@example.com"}
fmt.Println(user.Email) // Access field through pointer
```

### 4. **Error Handling**
```go
// Go uses explicit error handling
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to do something: %v", err)
}
```

### 5. **Goroutines and Channels**
```go
// Goroutines are lightweight threads
go func() {
    // This runs concurrently
}()

// Channels for communication between goroutines
ch := make(chan string)
go func() {
    ch <- "Hello from goroutine"
}()
msg := <-ch
```

## 🛠 Development Workflow

### 1. **Adding New Dependencies**
```bash
go get github.com/new-package
go mod tidy  # Clean up dependencies
```

### 2. **Running Tests**
```bash
go test ./...  # Run all tests
go test -v ./internal/auth  # Run tests with verbose output
```

### 3. **Building the Application**
```bash
go build -o fittrackplus cmd/server/main.go
./fittrackplus  # Run the binary
```

### 4. **Code Formatting**
```bash
go fmt ./...  # Format all Go files
go vet ./...  # Check for common mistakes
```

## 🔍 API Endpoints

### Health Check
- `GET /api/v1/health` - Check if the server is running

### Authentication (Coming Soon)
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Users (Coming Soon)
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile

## 🚧 Next Steps

1. **Implement Authentication Service**
   - JWT token generation and validation
   - Password hashing with bcrypt
   - Role-based access control

2. **Create User Management**
   - User registration and login
   - Profile management
   - Password reset functionality

3. **Build Plan Generation**
   - Fitness plan templates
   - Diet plan generation
   - Physiotherapy exercises

4. **Add Booking System**
   - Trainer availability
   - Session booking
   - Booking management

5. **Implement Progress Tracking**
   - Weight and measurement logging
   - Progress visualization
   - Analytics and reporting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License. 