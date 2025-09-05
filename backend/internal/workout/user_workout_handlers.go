package workout

import (
	"net/http"
	"strconv"

	"fittrackplus/internal/auth"
	"github.com/gin-gonic/gin"
)

// UserWorkoutHandler handles HTTP requests for user workout sessions
type UserWorkoutHandler struct {
	userWorkoutService *UserWorkoutService
}

// NewUserWorkoutHandler creates a new user workout handler
func NewUserWorkoutHandler(userWorkoutService *UserWorkoutService) *UserWorkoutHandler {
	return &UserWorkoutHandler{userWorkoutService: userWorkoutService}
}

// StartWorkout godoc
// @Summary Start a new workout session
// @Description Start a new workout session for the authenticated user
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body StartWorkoutRequest true "Workout start data"
// @Success 201 {object} UserWorkoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /workouts/start [post]
func (h *UserWorkoutHandler) StartWorkout(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req StartWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	userWorkout, err := h.userWorkoutService.StartWorkout(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to start workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Workout started successfully",
		"user_workout": userWorkout,
	})
}

// CompleteExerciseSet godoc
// @Summary Complete an exercise set
// @Description Record completion of an exercise set during a workout
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Param request body CompleteExerciseSetRequest true "Exercise set completion data"
// @Success 201 {object} UserWorkoutExerciseResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/{id}/complete-set [post]
func (h *UserWorkoutHandler) CompleteExerciseSet(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	var req CompleteExerciseSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	exerciseSet, err := h.userWorkoutService.CompleteExerciseSet(c.Request.Context(), uint(userWorkoutID), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to complete exercise set",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Exercise set completed successfully",
		"exercise_set": exerciseSet,
	})
}

// CompleteWorkout godoc
// @Summary Complete a workout session
// @Description Mark a workout session as completed
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Success 200 {object} UserWorkoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/{id}/complete [post]
func (h *UserWorkoutHandler) CompleteWorkout(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	userWorkout, err := h.userWorkoutService.CompleteWorkout(c.Request.Context(), uint(userWorkoutID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to complete workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout completed successfully",
		"user_workout": userWorkout,
	})
}

// PauseWorkout godoc
// @Summary Pause a workout session
// @Description Pause an active workout session
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Success 200 {object} UserWorkoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/{id}/pause [post]
func (h *UserWorkoutHandler) PauseWorkout(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	userWorkout, err := h.userWorkoutService.PauseWorkout(c.Request.Context(), uint(userWorkoutID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to pause workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout paused successfully",
		"user_workout": userWorkout,
	})
}

// ResumeWorkout godoc
// @Summary Resume a paused workout session
// @Description Resume a paused workout session
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Success 200 {object} UserWorkoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/{id}/resume [post]
func (h *UserWorkoutHandler) ResumeWorkout(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	userWorkout, err := h.userWorkoutService.ResumeWorkout(c.Request.Context(), uint(userWorkoutID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to resume workout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout resumed successfully",
		"user_workout": userWorkout,
	})
}

// GetUserWorkout godoc
// @Summary Get user workout session
// @Description Get detailed information about a specific user workout session
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Success 200 {object} UserWorkoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/session/{id} [get]
func (h *UserWorkoutHandler) GetUserWorkout(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	userWorkout, err := h.userWorkoutService.GetUserWorkoutByID(uint(userWorkoutID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout session not found"})
		return
	}

	// Verify user owns this workout session
	if userWorkout.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this workout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout session retrieved successfully",
		"user_workout": userWorkout,
	})
}

// GetUserWorkouts godoc
// @Summary Get user workout sessions
// @Description Get all workout sessions for the authenticated user
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (in_progress, completed, paused, cancelled)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /workouts/my-sessions [get]
func (h *UserWorkoutHandler) GetUserWorkouts(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	status := c.Query("status")

	userWorkouts, err := h.userWorkoutService.GetUserWorkouts(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch workout sessions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout sessions retrieved successfully",
		"user_workouts": userWorkouts,
	})
}

// GetWorkoutProgress godoc
// @Summary Get workout progress
// @Description Get the progress of a specific workout session
// @Tags user-workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Workout ID"
// @Success 200 {object} WorkoutProgressResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workouts/{id}/progress [get]
func (h *UserWorkoutHandler) GetWorkoutProgress(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userWorkoutIDStr := c.Param("id")
	userWorkoutID, err := strconv.ParseUint(userWorkoutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	progress, err := h.userWorkoutService.GetWorkoutProgress(uint(userWorkoutID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get workout progress",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workout progress retrieved successfully",
		"progress": progress,
	})
}
