package auth

import (
	"net/http"

	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: NewAuthService(cfg),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Register the user
	response, err := h.authService.Register(&req)
	if err != nil {
		// Check if it's a duplicate email error
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Login the user
	response, err := h.authService.Login(&req)
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "invalid email or password" || err.Error() == "account is deactivated" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to login",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the current authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get current user from context (set by middleware)
	user, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the current authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get current user ID from context
	userID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Update the user profile
	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
} 