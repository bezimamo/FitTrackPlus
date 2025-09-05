package exercise

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// ExerciseService handles exercise management business logic
type ExerciseService struct {
	db           *gorm.DB
	cfg          *config.Config
	cloudinary   *config.CloudinaryConfig
}

// NewExerciseService creates a new exercise service
func NewExerciseService(cfg *config.Config) *ExerciseService {
	return &ExerciseService{
		db:         database.GetDB(),
		cfg:        cfg,
		cloudinary: config.NewCloudinaryConfig(),
	}
}

// Request DTOs
type CreateExerciseRequest struct {
	Name         string                  `form:"name" binding:"required"`
	Description  string                  `form:"description"`
	Category     string                  `form:"category" binding:"required"` // strength, cardio, flexibility, balance
	MuscleGroup  string                  `form:"muscle_group" binding:"required"` // chest, back, legs, shoulders, arms, core, full_body
	Difficulty   string                  `form:"difficulty" binding:"required"` // beginner, intermediate, advanced
	Equipment    string                  `form:"equipment" binding:"required"` // bodyweight, dumbbells, barbell, machine, etc.
	Instructions string                  `form:"instructions"`
	Tips         string                  `form:"tips"`
	VideoFile    *multipart.FileHeader   `form:"video_file"` // Optional video file
	ImageFile    *multipart.FileHeader   `form:"image_file"` // Optional image file
}

type UpdateExerciseRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	Category     *string `json:"category,omitempty"`
	MuscleGroup  *string `json:"muscle_group,omitempty"`
	Difficulty   *string `json:"difficulty,omitempty"`
	Equipment    *string `json:"equipment,omitempty"`
	VideoURL     *string `json:"video_url,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	Instructions *string `json:"instructions,omitempty"`
	Tips         *string `json:"tips,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type ExerciseResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	MuscleGroup  string    `json:"muscle_group"`
	Difficulty   string    `json:"difficulty"`
	Equipment    string    `json:"equipment"`
	VideoURL     string    `json:"video_url"`
	ImageURL     string    `json:"image_url"`
	Instructions string    `json:"instructions"`
	Tips         string    `json:"tips"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateExercise creates a new exercise with optional media files (Admin only)
func (s *ExerciseService) CreateExercise(ctx context.Context, req *CreateExerciseRequest) (*ExerciseResponse, error) {
	// Check if exercise with same name already exists
	var existingExercise models.Exercise
	if err := s.db.Where("name = ?", req.Name).First(&existingExercise).Error; err == nil {
		return nil, errors.New("exercise with this name already exists")
	}

	// Initialize URLs
	videoURL := ""
	imageURL := ""

	// Upload video if provided
	if req.VideoFile != nil {
		if s.cloudinary == nil {
			return nil, errors.New("cloudinary not configured for video upload")
		}
		videoResult, err := s.cloudinary.UploadVideo(ctx, req.VideoFile, "exercises")
		if err != nil {
			return nil, fmt.Errorf("failed to upload video: %v", err)
		}
		videoURL = videoResult.SecureURL
	}

	// Upload image if provided
	if req.ImageFile != nil {
		if s.cloudinary == nil {
			return nil, errors.New("cloudinary not configured for image upload")
		}
		imageResult, err := s.cloudinary.UploadImage(ctx, req.ImageFile, "exercises")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %v", err)
		}
		imageURL = imageResult.SecureURL
	}

	// Create the exercise
	exercise := models.Exercise{
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		MuscleGroup:  req.MuscleGroup,
		Difficulty:   req.Difficulty,
		Equipment:    req.Equipment,
		VideoURL:     videoURL,
		ImageURL:     imageURL,
		Instructions: req.Instructions,
		Tips:         req.Tips,
		IsActive:     true,
	}

	if err := s.db.Create(&exercise).Error; err != nil {
		return nil, err
	}

	return s.buildExerciseResponse(&exercise), nil
}

// Media upload is now handled directly in CreateExercise endpoint

// GetExercises gets all exercises with optional filtering
func (s *ExerciseService) GetExercises(category, muscleGroup, difficulty, equipment, status string) ([]ExerciseResponse, error) {
	var exercises []models.Exercise
	query := s.db.Model(&models.Exercise{})

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if muscleGroup != "" {
		query = query.Where("muscle_group = ?", muscleGroup)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if equipment != "" {
		query = query.Where("equipment = ?", equipment)
	}
	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Order("name ASC").Find(&exercises).Error; err != nil {
		return nil, err
	}

	var responses []ExerciseResponse
	for _, exercise := range exercises {
		responses = append(responses, *s.buildExerciseResponse(&exercise))
	}

	return responses, nil
}

// GetExercise gets a specific exercise by ID
func (s *ExerciseService) GetExercise(exerciseID uint) (*ExerciseResponse, error) {
	var exercise models.Exercise
	if err := s.db.First(&exercise, exerciseID).Error; err != nil {
		return nil, errors.New("exercise not found")
	}

	return s.buildExerciseResponse(&exercise), nil
}

// UpdateExercise updates an existing exercise (Admin only)
func (s *ExerciseService) UpdateExercise(exerciseID uint, req *UpdateExerciseRequest) (*ExerciseResponse, error) {
	var exercise models.Exercise
	if err := s.db.First(&exercise, exerciseID).Error; err != nil {
		return nil, errors.New("exercise not found")
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if new name conflicts with existing exercise
		var existingExercise models.Exercise
		if err := s.db.Where("name = ? AND id != ?", *req.Name, exerciseID).First(&existingExercise).Error; err == nil {
			return nil, errors.New("exercise with this name already exists")
		}
		exercise.Name = *req.Name
	}
	if req.Description != nil {
		exercise.Description = *req.Description
	}
	if req.Category != nil {
		exercise.Category = *req.Category
	}
	if req.MuscleGroup != nil {
		exercise.MuscleGroup = *req.MuscleGroup
	}
	if req.Difficulty != nil {
		exercise.Difficulty = *req.Difficulty
	}
	if req.Equipment != nil {
		exercise.Equipment = *req.Equipment
	}
	if req.VideoURL != nil {
		exercise.VideoURL = *req.VideoURL
	}
	if req.ImageURL != nil {
		exercise.ImageURL = *req.ImageURL
	}
	if req.Instructions != nil {
		exercise.Instructions = *req.Instructions
	}
	if req.Tips != nil {
		exercise.Tips = *req.Tips
	}
	if req.IsActive != nil {
		exercise.IsActive = *req.IsActive
	}

	if err := s.db.Save(&exercise).Error; err != nil {
		return nil, err
	}

	return s.buildExerciseResponse(&exercise), nil
}

// DeleteExercise soft deletes an exercise (Admin only)
func (s *ExerciseService) DeleteExercise(exerciseID uint) error {
	var exercise models.Exercise
	if err := s.db.First(&exercise, exerciseID).Error; err != nil {
		return errors.New("exercise not found")
	}

	// Check if exercise is being used in any workouts
	var workoutCount int64
	if err := s.db.Model(&models.WorkoutExercise{}).Where("exercise_id = ?", exerciseID).Count(&workoutCount).Error; err != nil {
		return err
	}

	if workoutCount > 0 {
		return errors.New("cannot delete exercise - it is being used in workouts")
	}

	// Soft delete
	return s.db.Delete(&exercise).Error
}

// GetExerciseCategories gets all available exercise categories
func (s *ExerciseService) GetExerciseCategories() ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.Exercise{}).Distinct().Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetMuscleGroups gets all available muscle groups
func (s *ExerciseService) GetMuscleGroups() ([]string, error) {
	var muscleGroups []string
	if err := s.db.Model(&models.Exercise{}).Distinct().Pluck("muscle_group", &muscleGroups).Error; err != nil {
		return nil, err
	}
	return muscleGroups, nil
}

// GetDifficulties gets all available difficulty levels
func (s *ExerciseService) GetDifficulties() ([]string, error) {
	var difficulties []string
	if err := s.db.Model(&models.Exercise{}).Distinct().Pluck("difficulty", &difficulties).Error; err != nil {
		return nil, err
	}
	return difficulties, nil
}

// GetEquipment gets all available equipment types
func (s *ExerciseService) GetEquipment() ([]string, error) {
	var equipment []string
	if err := s.db.Model(&models.Exercise{}).Distinct().Pluck("equipment", &equipment).Error; err != nil {
		return nil, err
	}
	return equipment, nil
}

// GetExerciseStats gets exercise statistics for admin dashboard
func (s *ExerciseService) GetExerciseStats() (map[string]interface{}, error) {
	var totalExercises, activeExercises, strengthExercises, cardioExercises, flexibilityExercises int64
	var avgDifficulty string

	// Get total exercises
	if err := s.db.Model(&models.Exercise{}).Count(&totalExercises).Error; err != nil {
		return nil, err
	}

	// Get active exercises
	if err := s.db.Model(&models.Exercise{}).Where("is_active = ?", true).Count(&activeExercises).Error; err != nil {
		return nil, err
	}

	// Get exercises by category
	if err := s.db.Model(&models.Exercise{}).Where("category = ?", "strength").Count(&strengthExercises).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Exercise{}).Where("category = ?", "cardio").Count(&cardioExercises).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Exercise{}).Where("category = ?", "flexibility").Count(&flexibilityExercises).Error; err != nil {
		return nil, err
	}

	// Calculate average difficulty (simplified - most common difficulty)
	var difficulties []string
	if err := s.db.Model(&models.Exercise{}).Pluck("difficulty", &difficulties).Error; err != nil {
		return nil, err
	}

	// Count occurrences of each difficulty
	difficultyCount := make(map[string]int)
	for _, diff := range difficulties {
		difficultyCount[diff]++
	}

	// Find most common difficulty
	maxCount := 0
	for diff, count := range difficultyCount {
		if count > maxCount {
			maxCount = count
			avgDifficulty = diff
		}
	}

	return map[string]interface{}{
		"total_exercises":      totalExercises,
		"active_exercises":     activeExercises,
		"strength_exercises":   strengthExercises,
		"cardio_exercises":     cardioExercises,
		"flexibility_exercises": flexibilityExercises,
		"avg_difficulty":       avgDifficulty,
	}, nil
}

// Helper function to build response
func (s *ExerciseService) buildExerciseResponse(exercise *models.Exercise) *ExerciseResponse {
	return &ExerciseResponse{
		ID:           exercise.ID,
		Name:         exercise.Name,
		Description:  exercise.Description,
		Category:     exercise.Category,
		MuscleGroup:  exercise.MuscleGroup,
		Difficulty:   exercise.Difficulty,
		Equipment:    exercise.Equipment,
		VideoURL:     exercise.VideoURL,
		ImageURL:     exercise.ImageURL,
		Instructions: exercise.Instructions,
		Tips:         exercise.Tips,
		IsActive:     exercise.IsActive,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}
}