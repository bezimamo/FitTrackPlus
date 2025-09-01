package workout

import (
	"net/http"
	"strconv"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// WorkoutHandler handles workout HTTP requests
type WorkoutHandler struct {
	workoutService *WorkoutService
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(cfg *config.Config) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: NewWorkoutService(cfg),
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
	workout, err := h.workoutService.CreateWorkout(trainerID, &req)
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

	// Get query parameters
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
	workouts, err := h.workoutService.GetTrainerWorkouts(trainerID, memberID, status)
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

	// Get member ID
	memberID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get query parameters
	status := c.Query("status")

	// Get workouts
	workouts, err := h.workoutService.GetMemberWorkouts(memberID, status)
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

	// Get user info
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Get workout
	workout, err := h.workoutService.GetWorkout(uint(workoutID), userID, userRole)
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
	workout, err := h.workoutService.UpdateWorkout(uint(workoutID), trainerID, &req)
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
