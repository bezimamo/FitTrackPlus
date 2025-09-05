package workout

import (
	"net/http"
	"strconv"

	"fittrackplus/internal/auth"

	"github.com/gin-gonic/gin"
)

// WorkoutHandler handles workout HTTP requests
type WorkoutHandler struct {
	workoutService *WorkoutService
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(workoutService *WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: workoutService,
	}
}

// CreateWorkout godoc
// @Summary Create a new workout for a member (Trainer only)
// @Description Create a personalized workout session for an assigned member
// @Tags Workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workout body CreateWorkoutRequest true "Workout details"
// @Success 201 {object} WorkoutResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer only"
// @Router /workouts [post]
func (h *WorkoutHandler) CreateWorkout(c *gin.Context) {
	// Check if user is trainer
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can create workouts"})
		return
	}

	// Get trainer ID
	trainerID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Bind request
	var req CreateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create workout
	workout, err := h.workoutService.CreateWorkout(c.Request.Context(), &req, trainerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Workout created successfully!",
		"workout": workout,
	})
}

// GetTrainerWorkouts godoc
// @Summary Get workouts created by trainer (Trainer only)
// @Description Get all workouts created by the current trainer with optional filtering
// @Tags Workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param member_id query int false "Filter by specific member"
// @Param status query string false "Filter by status (active, inactive)"
// @Success 200 {array} WorkoutResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer only"
// @Router /trainer/workouts [get]
func (h *WorkoutHandler) GetTrainerWorkouts(c *gin.Context) {
	// Check if user is trainer
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can view their workouts"})
		return
	}

	// Get trainer ID
	trainerID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	memberIDStr := c.Query("member_id")
	status := c.Query("status")

	var memberID *uint
	if memberIDStr != "" {
		if id, err := strconv.ParseUint(memberIDStr, 10, 32); err == nil {
			memberIDUint := uint(id)
			memberID = &memberIDUint
		}
	}

	// Get workouts
	workouts, err := h.workoutService.GetWorkouts("", "", "", &trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get workouts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your workouts",
		"workouts": workouts,
		"total": len(workouts),
		"filters": gin.H{
			"member_id": memberID,
			"status": status,
		},
	})
}

// GetMemberWorkouts godoc
// @Summary Get workouts assigned to member (Member only)
// @Description Get all workouts assigned to the current member
// @Tags Workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (active, inactive)"
// @Success 200 {array} WorkoutResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Member only"
// @Router /workouts/my-workouts [get]
func (h *WorkoutHandler) GetMemberWorkouts(c *gin.Context) {
	// Check if user is member
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "member" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only members can view their workouts"})
		return
	}

	// Get member ID (for future use if needed)
	_, exists = auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get query parameters
	status := c.Query("status")

	// Get workouts
	workouts, err := h.workoutService.GetWorkouts("", "", "", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get workouts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your assigned workouts",
		"workouts": workouts,
		"total": len(workouts),
		"filter": status,
	})
}

// GetWorkout godoc
// @Summary Get specific workout details
// @Description Get detailed information about a specific workout
// @Tags Workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workout ID"
// @Success 200 {object} WorkoutResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Access denied"
// @Router /workouts/{id} [get]
func (h *WorkoutHandler) GetWorkout(c *gin.Context) {
	// Get workout ID from URL
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID"})
		return
	}

	// Get user info (for future use if needed)
	_, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	_, exists = auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Get workout
	workout, err := h.workoutService.GetWorkoutByID(uint(workoutID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout details",
		"workout": workout,
	})
}

// UpdateWorkout godoc
// @Summary Update an existing workout (Trainer only)
// @Description Update workout details and exercises
// @Tags Workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workout ID"
// @Param workout body CreateWorkoutRequest true "Updated workout details"
// @Success 200 {object} WorkoutResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer only"
// @Router /workouts/{id} [put]
func (h *WorkoutHandler) UpdateWorkout(c *gin.Context) {
	// Check if user is trainer
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can update workouts"})
		return
	}

	// Get trainer ID
	trainerID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get workout ID from URL
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID"})
		return
	}

	// Bind request
	var req CreateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update workout
	workout, err := h.workoutService.UpdateWorkout(c.Request.Context(), uint(workoutID), &req, trainerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout updated successfully!",
		"workout": workout,
	})
}

// GetWorkouts godoc
// @Summary Get all workouts (public access)
// @Description Get all active workout templates with optional filtering
// @Tags Workouts
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param difficulty query string false "Filter by difficulty"
// @Param search query string false "Search by name or description"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workouts [get]
func (h *WorkoutHandler) GetWorkouts(c *gin.Context) {
	// Get query parameters
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	search := c.Query("search")

	// Get workouts
	workouts, err := h.workoutService.GetWorkouts(category, difficulty, search, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get workouts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workouts retrieved successfully",
		"workouts": workouts,
	})
}

// GetWorkoutByID godoc
// @Summary Get workout by ID (public access)
// @Description Get a specific workout template by ID
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/workouts/{id} [get]
func (h *WorkoutHandler) GetWorkoutByID(c *gin.Context) {
	// Get workout ID from URL
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID"})
		return
	}

	// Get workout
	workout, err := h.workoutService.GetWorkoutByID(uint(workoutID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout retrieved successfully",
		"workout": workout,
	})
}

// GetWorkoutCategories godoc
// @Summary Get workout categories (public access)
// @Description Get all available workout categories
// @Tags Workouts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workouts/categories [get]
func (h *WorkoutHandler) GetWorkoutCategories(c *gin.Context) {
	categories, err := h.workoutService.GetWorkoutCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"categories": categories,
	})
}

// GetWorkoutDifficulties godoc
// @Summary Get workout difficulties (public access)
// @Description Get all available workout difficulty levels
// @Tags Workouts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workouts/difficulties [get]
func (h *WorkoutHandler) GetWorkoutDifficulties(c *gin.Context) {
	difficulties, err := h.workoutService.GetWorkoutDifficulties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get difficulties",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Difficulties retrieved successfully",
		"difficulties": difficulties,
	})
}

// DeleteWorkout godoc
// @Summary Delete a workout (Admin/Trainer only)
// @Description Soft delete a workout template
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/admin/workouts/{id} [delete]
func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	// Check if user is trainer or admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers and admins can delete workouts"})
		return
	}

	// Get user ID
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get workout ID from URL
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID"})
		return
	}

	// Delete workout
	err = h.workoutService.DeleteWorkout(uint(workoutID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout deleted successfully",
	})
}

// GetWorkoutStats godoc
// @Summary Get workout statistics (Admin only)
// @Description Get comprehensive workout statistics
// @Tags Workouts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/workouts/stats [get]
func (h *WorkoutHandler) GetWorkoutStats(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view workout statistics"})
		return
	}

	// Get stats
	stats, err := h.workoutService.GetWorkoutStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get workout statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout statistics retrieved successfully",
		"stats": stats,
	})
}
