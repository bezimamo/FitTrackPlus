package auth

import (
	"testing"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"
)

func setupTestDB(t *testing.T) *config.Config {
	// Load test configuration
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "fiker2901",
		DBName:     "fittrackplus_test", // Use test database
		JWTSecret:  "test-secret-key",
	}

	// Connect to test database
	err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up test data
	db := database.GetDB()
	db.Exec("DELETE FROM users")

	return cfg
}

func TestAuthService_Register(t *testing.T) {
	cfg := setupTestDB(t)
	authService := NewAuthService(cfg)

	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
	}{
		{
			name: "Valid registration",
			req: &RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+1234567890",
				Role:      "member",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			req: &RegisterRequest{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "member",
			},
			wantErr: true,
		},
		{
			name: "Short password",
			req: &RegisterRequest{
				Email:     "test2@example.com",
				Password:  "123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "member",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Register(tt.req)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			if response.Token == "" {
				t.Errorf("Expected token but got empty string")
			}

			if response.User.Email != tt.req.Email {
				t.Errorf("Expected email %s, got %s", tt.req.Email, response.User.Email)
			}

			if response.User.Password != "" {
				t.Errorf("Password should not be returned in response")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	cfg := setupTestDB(t)
	authService := NewAuthService(cfg)

	// First register a user
	registerReq := &RegisterRequest{
		Email:     "login@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "member",
	}

	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user for login test: %v", err)
	}

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
	}{
		{
			name: "Valid login",
			req: &LoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Invalid password",
			req: &LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "Non-existent user",
			req: &LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.req)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			if response.Token == "" {
				t.Errorf("Expected token but got empty string")
			}

			if response.User.Email != tt.req.Email {
				t.Errorf("Expected email %s, got %s", tt.req.Email, response.User.Email)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	cfg := setupTestDB(t)
	authService := NewAuthService(cfg)

	// Register and login to get a token
	registerReq := &RegisterRequest{
		Email:     "token@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "member",
	}

	registerResp, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user for token test: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantErr    bool
		wantUserID uint
	}{
		{
			name:       "Valid token",
			token:      registerResp.Token,
			wantErr:    false,
			wantUserID: registerResp.User.ID,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Errorf("Expected claims but got nil")
				return
			}

			if claims.UserID != tt.wantUserID {
				t.Errorf("Expected user ID %d, got %d", tt.wantUserID, claims.UserID)
			}

			if claims.Email != registerReq.Email {
				t.Errorf("Expected email %s, got %s", registerReq.Email, claims.Email)
			}
		})
	}
} 