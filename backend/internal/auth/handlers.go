package auth

import (
	"net/http"
	"strconv"
	"time"

	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Phone     string `json:"phone" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=member trainer physio admin"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone     *string `json:"phone,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// ChangeRoleRequest represents role change request
type ChangeRoleRequest struct {
	NewRole string `json:"new_role" binding:"required,oneof=member trainer physio admin"`
}

// UserListResponse represents user list response
type UserListResponse struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



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
// @Success 200 {object} map[string]interface{}
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
// @Success 200 {object} map[string]interface{}
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

// GetTrainers godoc
// @Summary Get all trainers (Admin only)
// @Description Get a list of all users with trainer role
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} TrainerInfo
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/users/trainers [get]
func (h *AuthHandler) GetTrainers(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only admins can view all trainers",
		})
		return
	}

	// Get trainers
	trainers, err := h.authService.GetTrainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get trainers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainers retrieved successfully",
		"trainers": trainers,
		"total": len(trainers),
	})
}

// GetAllUsers godoc
// @Summary Get all users with pagination and filtering (Admin only)
// @Description Get a paginated list of all users with optional role, search, and status filters
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param role query string false "Filter by role"
// @Param search query string false "Search in name or email"
// @Param status query string false "Filter by status (active/inactive)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/users [get]
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can view all users",
		})
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")
	search := c.Query("search")
	status := c.Query("status")

	users, total, err := h.authService.GetAllUsers(page, limit, role, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUserByID godoc
// @Summary Get user by ID (Admin only)
// @Description Get detailed information about a specific user
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can view user details",
		})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.authService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"user": user,
	})
}

// CreateUser godoc
// @Summary Create new user (Admin only)
// @Description Create a new user with specified role and details
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can create users",
		})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": user,
	})
}

// UpdateUser godoc
// @Summary Update user (Admin only)
// @Description Update an existing user's information
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/users/{id} [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can update users",
		})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.UpdateUser(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": user,
	})
}

// DeactivateUser godoc
// @Summary Deactivate user (Admin only)
// @Description Deactivate a user account (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/users/{id} [delete]
func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can deactivate users",
		})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.authService.DeactivateUser(uint(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to deactivate user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
	})
}

// GetAvailableRoles godoc
// @Summary Get available user roles (Admin only)
// @Description Get a list of all available user roles
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/users/roles [get]
func (h *AuthHandler) GetAvailableRoles(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can view roles",
		})
		return
	}

	roles := h.authService.GetAvailableRoles()
	c.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"roles": roles,
	})
}

// ChangeUserRole godoc
// @Summary Change user role (Admin only)
// @Description Change a user's role to a new role
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body ChangeRoleRequest true "Role change data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/users/{id}/change-role [post]
func (h *AuthHandler) ChangeUserRole(c *gin.Context) {
	// Check if user is admin
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can change user roles",
		})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.ChangeUserRole(uint(userID), req.NewRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to change user role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role changed successfully",
		"user": user,
	})
}

