package profile

import (
	"net/http"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles HTTP requests for profile management
type ProfileHandler struct {
	profileService     *ProfileService
	roleProfileService *RoleProfileService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		profileService:     NewProfileService(cfg),
		roleProfileService: NewRoleProfileService(cfg),
	}
}

// SetupProfile handles complete profile setup
// @Summary Setup user profile
// @Description Complete profile setup for new users
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProfileSetupRequest true "Profile setup data"
// @Success 200 {object} ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/setup [post]
func (h *ProfileHandler) SetupProfile(c *gin.Context) {
	var req ProfileSetupRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get current user ID from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Setup the profile
	response, err := h.profileService.SetupProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to setup profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SetupRoleProfile handles role-based profile setup
// @Summary Setup role-based profile
// @Description Complete profile setup based on user role (member, trainer, admin, physio)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RoleProfileSetupRequest true "Role-based profile setup data"
// @Success 200 {object} RoleProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/setup-role [post]
func (h *ProfileHandler) SetupRoleProfile(c *gin.Context) {
	var req RoleProfileSetupRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get current user info from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	// Setup the role-based profile
	response, err := h.roleProfileService.SetupRoleProfile(userID, userRole, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to setup role-based profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRoleProfile handles role-based profile retrieval
// @Summary Get role-based profile
// @Description Get the current user's role-specific profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RoleProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/role [get]
func (h *ProfileHandler) GetRoleProfile(c *gin.Context) {
	// Get current user info from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	// Get the role-based profile
	response, err := h.roleProfileService.GetRoleProfile(userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get role-based profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CheckRoleProfileCompletion handles role-based profile completion check
// @Summary Check role-based profile completion
// @Description Check profile completion percentage based on user role
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RoleProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/role/completion [get]
func (h *ProfileHandler) CheckRoleProfileCompletion(c *gin.Context) {
	// Get current user info from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	// Check the role-based profile completion
	response, err := h.roleProfileService.CheckProfileCompletion(userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check role-based profile completion",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile handles profile retrieval
// @Summary Get user profile
// @Description Get the current user's complete profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get current user ID from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get the profile
	response, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UploadProfileImage handles profile image upload
// @Summary Upload profile image
// @Description Upload a profile image for the current user
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Profile image file"
// @Success 200 {object} ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/upload-image [post]
func (h *ProfileHandler) UploadProfileImage(c *gin.Context) {
	// Get current user ID from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
			"details": err.Error(),
		})
		return
	}

	// Upload the image
	response, err := h.profileService.UploadProfileImage(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upload image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CheckProfileCompletion checks if profile is complete
// @Summary Check profile completion
// @Description Check if the current user's profile is complete
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/profile/completion [get]
func (h *ProfileHandler) CheckProfileCompletion(c *gin.Context) {
	// Get current user ID from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get the profile
	response, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_complete": response.IsComplete,
		"completion_percentage": response.Completion,
		"profile_exists": response.UserProfile.ID > 0,
	})
} 