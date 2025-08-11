package auth

import (
	"errors"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTClaims represents the claims in our JWT token
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles all authentication-related operations
type AuthService struct {
	db       *gorm.DB
	jwtSecret string
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		db:        database.GetDB(),
		jwtSecret: cfg.JWTSecret,
	}
}

// RegisterRequest represents the data needed for user registration
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Role      string `json:"role" binding:"required,oneof=member trainer physio admin"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the data needed for profile updates
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Register creates a new user account
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		IsActive:  true,
	}

	// Save user to database
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate JWT token
	token, expiresAt, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	// Don't return the password in the response
	user.Password = ""

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	// Don't return the password in the response
	user.Password = ""

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the user claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	
	// Don't return the password
	user.Password = ""
	return &user, nil
}

// UpdateProfile updates a user's profile information
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Update the user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	// Save the updated user
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return &user, nil
}

// generateToken creates a new JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, time.Time, error) {
	// Set expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fittrackplus",
			Subject:   user.Email,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
} 