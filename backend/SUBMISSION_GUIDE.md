# FitTrack+ Backend - Student Project Submission

## 🎯 **Project Overview**
**FitTrack+** is a comprehensive fitness platform backend built with Go, featuring user authentication, profile management, and a complete API system.

## ✅ **Completed Features**

### **1. Authentication System**
- ✅ User Registration (`POST /api/v1/auth/register`)
- ✅ User Login (`POST /api/v1/auth/login`)
- ✅ JWT Token Authentication
- ✅ Password Hashing (bcrypt)
- ✅ Protected Routes

### **2. Enhanced Profile Management**
- ✅ Complete Profile Setup (`POST /api/v1/users/profile/setup`)
- ✅ Profile Retrieval (`GET /api/v1/users/profile`)
- ✅ Profile Image Upload (`POST /api/v1/users/profile/upload-image`)
- ✅ Profile Completion Tracking (`GET /api/v1/users/profile/completion`)
- ✅ Comprehensive Profile Fields (fitness goals, medical info, measurements)

### **3. Database System**
- ✅ PostgreSQL Integration
- ✅ GORM ORM
- ✅ Auto-migration
- ✅ Complete Data Models

### **4. API Documentation**
- ✅ Swagger/OpenAPI Documentation
- ✅ Interactive API Testing
- ✅ Complete Endpoint Documentation

## 🛠 **Technology Stack**
- **Language**: Go (Golang)
- **Framework**: Gin (HTTP Web Framework)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Password Hashing**: bcrypt

## 📁 **Project Structure**
```
backend/
├── cmd/server/main.go          # Main application entry point
├── internal/
│   ├── auth/                   # Authentication system
│   │   ├── auth.go            # Business logic
│   │   ├── handlers.go        # HTTP handlers
│   │   └── middleware.go      # JWT middleware
│   ├── profile/               # Profile management
│   │   ├── profile.go         # Business logic
│   │   └── handlers.go        # HTTP handlers
│   └── common/
│       ├── config/            # Configuration management
│       ├── database/          # Database connection
│       └── models/            # Data models
├── docs/                      # Swagger documentation
├── migrations/                # Database migrations
├── scripts/                   # Setup scripts
├── go.mod                     # Go dependencies
├── go.sum                     # Dependency checksums
├── .env                       # Environment variables
└── README.md                  # Project documentation
```

## 🚀 **Quick Start Guide**

### **Prerequisites**
1. **Go** (version 1.21 or higher)
2. **PostgreSQL** (version 12 or higher)
3. **Git**

### **Installation Steps**

#### **Step 1: Clone and Setup**
```bash
# Clone the project
git clone <repository-url>
cd FitTrackPlus/backend

# Install Go dependencies
go mod tidy
```

#### **Step 2: Database Setup**
```bash
# Install PostgreSQL (if not already installed)
# Windows: Download from https://www.postgresql.org/download/windows/
# macOS: brew install postgresql
# Linux: sudo apt-get install postgresql

# Create database
psql -U postgres
CREATE DATABASE fittrackplus;
CREATE USER fittrackplus_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE fittrackplus TO fittrackplus_user;
\q
```

#### **Step 3: Environment Configuration**
```bash
# Copy environment template
cp env.example .env

# Edit .env file with your database credentials
DB_HOST=localhost
DB_PORT=5432
DB_USER=fittrackplus_user
DB_PASSWORD=your_password
DB_NAME=fittrackplus
JWT_SECRET=your_super_secret_jwt_key_here
PORT=8080
```

#### **Step 4: Run the Application**
```bash
# Start the server
go run cmd/server/main.go
```

#### **Step 5: Access the Application**
- **API Documentation**: http://localhost:8080/api/docs
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/v1/health

## 📋 **API Endpoints**

### **Authentication (Public)**
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### **User Profile (Protected - Requires JWT)**
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update basic profile
- `POST /api/v1/users/profile/setup` - Complete profile setup
- `POST /api/v1/users/profile/upload-image` - Upload profile image
- `GET /api/v1/users/profile/completion` - Check profile completion

### **System**
- `GET /api/v1/health` - Health check

## 🧪 **Testing the API**

### **1. Register a New User**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "member"
  }'
```

### **2. Login**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### **3. Setup Complete Profile (Use token from login)**
```bash
curl -X POST http://localhost:8080/api/v1/users/profile/setup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "height": 175.0,
    "weight": 70.0,
    "age": 25,
    "gender": "male",
    "goals": ["lose_weight", "gain_muscle"],
    "target_weight": 65.0,
    "timeline": 90,
    "medical_history": "No major issues",
    "preferred_workout_time": "morning",
    "workout_days": ["monday", "wednesday", "friday"],
    "communication_preference": "email"
  }'
```

## 📊 **Database Schema**

### **Users Table**
- `id` (Primary Key)
- `email` (Unique)
- `password` (Hashed)
- `first_name`, `last_name`
- `role` (member, trainer, physio, admin)
- `phone`
- `is_active`
- `created_at`, `updated_at`

### **User Profiles Table**
- `id` (Primary Key)
- `user_id` (Foreign Key)
- `height`, `weight`, `age`, `gender`
- `goals` (JSON)
- `target_weight`, `timeline`
- `medical_history`, `allergies`, `medications`
- `body_fat_percentage`, `muscle_mass`
- `body_measurements` (JSON)
- `profile_image_url`
- `preferred_workout_time`, `workout_days` (JSON)
- `communication_preference`
- `is_profile_complete`
- `created_at`, `updated_at`

## 🔒 **Security Features**
- ✅ Password hashing with bcrypt
- ✅ JWT token authentication
- ✅ Protected API endpoints
- ✅ Input validation
- ✅ CORS configuration
- ✅ Environment variable configuration

## 📈 **Learning Outcomes**
This project demonstrates:
1. **Go Language Proficiency**: Modern Go development practices
2. **Web API Development**: RESTful API design and implementation
3. **Database Design**: PostgreSQL with GORM ORM
4. **Authentication**: JWT-based authentication system
5. **Security**: Password hashing, input validation
6. **Documentation**: API documentation with Swagger
7. **Project Structure**: Clean architecture and code organization

## 🎓 **Student Information**
- **Project**: FitTrack+ Backend API
- **Technology**: Go, PostgreSQL, JWT Authentication
- **Features**: Complete authentication and profile management system
- **Status**: MVP Ready (Authentication + Profile Management)

## 📞 **Support**
For any questions about this project:
- Check the README.md file
- Review the Swagger documentation
- Test the API endpoints using the provided examples

---
**Note**: This is a student project demonstrating backend development skills with Go and PostgreSQL. 