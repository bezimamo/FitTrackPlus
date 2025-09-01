package exercise

import (
	"net/http"
	"strconv"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// ExerciseHandler handles exercise HTTP requests
type ExerciseHandler struct {
	exerciseService *ExerciseService
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(cfg *config.Config) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: NewExerciseService(cfg),
	}
}

// CreateExercise godoc
// @Summary Create a new exercise with optional media files (Admin only)
// @Description Add a new exercise to the exercise library with optional video and image uploads
// @Tags Exercises
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Exercise name"
// @Param description formData string false "Exercise description"
// @Param category formData string true "Exercise category (strength, cardio, flexibility, balance)"
// @Param muscle_group formData string true "Target muscle group"
// @Param difficulty formData string true "Difficulty level (beginner, intermediate, advanced)"
// @Param equipment formData string true "Required equipment"
// @Param instructions formData string false "Exercise instructions"
// @Param tips formData string false "Exercise tips"
// @Param video_file formData file false "Video file (MP4, AVI, MOV, etc.)"
// @Param image_file formData file false "Image file (JPG, PNG, etc.)"
// @Success 201 {object} ExerciseResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/exercises [post]
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create exercises"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse form data",
			"details": err.Error(),
		})
		return
	}

	// Bind request
	var req CreateExerciseRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create exercise with media files
	exercise, err := h.exerciseService.CreateExercise(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create exercise",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Exercise created successfully with media!",
		"exercise": exercise,
	})
}

// GetExercises godoc
// @Summary Get all exercises with optional filtering
// @Description Retrieve exercises with optional filtering by category, muscle group, difficulty, equipment, and status
// @Tags Exercises
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param muscle_group query string false "Filter by muscle group"
// @Param difficulty query string false "Filter by difficulty"
// @Param equipment query string false "Filter by equipment"
// @Param status query string false "Filter by status (active, inactive)"
// @Success 200 {array} ExerciseResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /exercises [get]
func (h *ExerciseHandler) GetExercises(c *gin.Context) {
	// Get query parameters
	category := c.Query("category")
	muscleGroup := c.Query("muscle_group")
	difficulty := c.Query("difficulty")
	equipment := c.Query("equipment")
	status := c.Query("status")

	// Get exercises
	exercises, err := h.exerciseService.GetExercises(category, muscleGroup, difficulty, equipment, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get exercises",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercises retrieved successfully",
		"exercises": exercises,
		"total": len(exercises),
		"filters": gin.H{
			"category": category,
			"muscle_group": muscleGroup,
			"difficulty": difficulty,
			"equipment": equipment,
			"status": status,
		},
	})
}

// GetExercise godoc
// @Summary Get specific exercise details
// @Description Get detailed information about a specific exercise
// @Tags Exercises
// @Accept json
// @Produce json
// @Param id path int true "Exercise ID"
// @Success 200 {object} ExerciseResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /exercises/{id} [get]
func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	// Get exercise ID from URL
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	// Get exercise
	exercise, err := h.exerciseService.GetExercise(uint(exerciseID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get exercise",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise details",
		"exercise": exercise,
	})
}

// UpdateExercise godoc
// @Summary Update an existing exercise (Admin only)
// @Description Update exercise details
// @Tags Exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exercise ID"
// @Param exercise body UpdateExerciseRequest true "Updated exercise details"
// @Success 200 {object} ExerciseResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/exercises/{id} [put]
func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update exercises"})
		return
	}

	// Get exercise ID from URL
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	// Bind request
	var req UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update exercise
	exercise, err := h.exerciseService.UpdateExercise(uint(exerciseID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update exercise",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise updated successfully!",
		"exercise": exercise,
	})
}

// DeleteExercise godoc
// @Summary Delete an exercise (Admin only)
// @Description Soft delete an exercise (cannot delete if used in workouts)
// @Tags Exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exercise ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/exercises/{id} [delete]
func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete exercises"})
		return
	}

	// Get exercise ID from URL
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	// Delete exercise
	if err := h.exerciseService.DeleteExercise(uint(exerciseID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete exercise",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise deleted successfully!",
	})
}

// Media upload is now handled directly in CreateExercise endpoint

// GetExerciseCategories godoc
// @Summary Get all available exercise categories
// @Description Retrieve all unique exercise categories for filtering
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {array} string "List of categories"
// @Router /exercises/categories [get]
func (h *ExerciseHandler) GetExerciseCategories(c *gin.Context) {
	categories, err := h.exerciseService.GetExerciseCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise categories",
		"categories": categories,
	})
}

// GetMuscleGroups godoc
// @Summary Get all available muscle groups
// @Description Retrieve all unique muscle groups for filtering
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {array} string "List of muscle groups"
// @Router /exercises/muscle-groups [get]
func (h *ExerciseHandler) GetMuscleGroups(c *gin.Context) {
	muscleGroups, err := h.exerciseService.GetMuscleGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get muscle groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Muscle groups",
		"muscle_groups": muscleGroups,
	})
}

// GetDifficulties godoc
// @Summary Get all available difficulty levels
// @Description Retrieve all unique difficulty levels for filtering
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {array} string "List of difficulties"
// @Router /exercises/difficulties [get]
func (h *ExerciseHandler) GetDifficulties(c *gin.Context) {
	difficulties, err := h.exerciseService.GetDifficulties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get difficulties",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Difficulty levels",
		"difficulties": difficulties,
	})
}

// GetEquipment godoc
// @Summary Get all available equipment types
// @Description Retrieve all unique equipment types for filtering
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {array} string "List of equipment"
// @Router /exercises/equipment [get]
func (h *ExerciseHandler) GetEquipment(c *gin.Context) {
	equipment, err := h.exerciseService.GetEquipment()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get equipment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Equipment types",
		"equipment": equipment,
	})
}
