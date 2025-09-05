package workout

import (
	"context"
	"errors"
	"fmt"

	"fittrackplus/internal/common/models"
	"gorm.io/gorm"
)

// WorkoutService handles workout business logic
type WorkoutService struct {
	db *gorm.DB
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(db *gorm.DB) *WorkoutService {
	return &WorkoutService{db: db}
}

// CreateWorkoutRequest represents the request to create a workout
type CreateWorkoutRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Description string                    `json:"description"`
	Category    string                    `json:"category" binding:"required"`
	Difficulty  string                    `json:"difficulty" binding:"required"`
	Duration    int                       `json:"duration"`
	Equipment   string                    `json:"equipment"`
	IsTemplate  bool                      `json:"is_template"`
	Exercises   []CreateWorkoutExerciseRequest `json:"exercises"`
}

// CreateWorkoutExerciseRequest represents an exercise in a workout
type CreateWorkoutExerciseRequest struct {
	ExerciseID uint     `json:"exercise_id" binding:"required"`
	Order      int      `json:"order" binding:"required"`
	Sets       int      `json:"sets" binding:"required"`
	Reps       int      `json:"reps"`
	Weight     float64  `json:"weight"`
	Duration   int      `json:"duration"`
	RestTime   int      `json:"rest_time"`
	Notes      string   `json:"notes"`
}

// WorkoutResponse represents the response for a workout
type WorkoutResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Category    string                  `json:"category"`
	Difficulty  string                  `json:"difficulty"`
	Duration    int                     `json:"duration"`
	Equipment   string                  `json:"equipment"`
	IsTemplate  bool                    `json:"is_template"`
	IsActive    bool                    `json:"is_active"`
	CreatedBy   uint                    `json:"created_by"`
	CreatedAt   string                  `json:"created_at"`
	UpdatedAt   string                  `json:"updated_at"`
	Exercises   []WorkoutExerciseResponse `json:"exercises,omitempty"`
	CreatedByUser *UserInfo             `json:"created_by_user,omitempty"`
}

// WorkoutExerciseResponse represents an exercise in a workout response
type WorkoutExerciseResponse struct {
	ID         uint           `json:"id"`
	ExerciseID uint           `json:"exercise_id"`
	Order      int            `json:"order"`
	Sets       int            `json:"sets"`
	Reps       int            `json:"reps"`
	Weight     float64        `json:"weight"`
	Duration   int            `json:"duration"`
	RestTime   int            `json:"rest_time"`
	Notes      string         `json:"notes"`
	Exercise   *ExerciseInfo  `json:"exercise,omitempty"`
}

// ExerciseInfo represents basic exercise information
type ExerciseInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MuscleGroup string `json:"muscle_group"`
	Difficulty  string `json:"difficulty"`
	Equipment   string `json:"equipment"`
	ImageURL    string `json:"image_url,omitempty"`
	VideoURL    string `json:"video_url,omitempty"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// WorkoutStats represents workout statistics
type WorkoutStats struct {
	TotalWorkouts     int64 `json:"total_workouts"`
	ActiveWorkouts    int64 `json:"active_workouts"`
	StrengthWorkouts  int64 `json:"strength_workouts"`
	CardioWorkouts    int64 `json:"cardio_workouts"`
	HIITWorkouts      int64 `json:"hiit_workouts"`
	FlexibilityWorkouts int64 `json:"flexibility_workouts"`
	AvgDuration       int64 `json:"avg_duration"`
}

// CreateWorkout creates a new workout
func (s *WorkoutService) CreateWorkout(ctx context.Context, req *CreateWorkoutRequest, createdBy uint) (*WorkoutResponse, error) {
	// Validate category
	validCategories := []string{"strength", "cardio", "hiit", "flexibility", "yoga", "pilates", "crossfit"}
	categoryValid := false
	for _, cat := range validCategories {
		if cat == req.Category {
			categoryValid = true
			break
		}
	}
	if !categoryValid {
		return nil, errors.New("invalid category: must be one of strength, cardio, hiit, flexibility, yoga, pilates, crossfit")
	}

	// Validate difficulty
	validDifficulties := []string{"beginner", "intermediate", "advanced"}
	difficultyValid := false
	for _, diff := range validDifficulties {
		if diff == req.Difficulty {
			difficultyValid = true
			break
		}
	}
	if !difficultyValid {
		return nil, errors.New("invalid difficulty: must be one of beginner, intermediate, advanced")
	}

	// Check if workout with same name already exists
	var existingWorkout models.Workout
	if err := s.db.Where("name = ? AND created_by = ?", req.Name, createdBy).First(&existingWorkout).Error; err == nil {
		return nil, errors.New("workout with this name already exists")
	}

	// Create the workout
	workout := models.Workout{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Difficulty:  req.Difficulty,
		Duration:    req.Duration,
		Equipment:   req.Equipment,
		IsTemplate:  req.IsTemplate,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if err := s.db.Create(&workout).Error; err != nil {
		return nil, err
	}

	// Add exercises to the workout
	for _, exerciseReq := range req.Exercises {
		// Verify exercise exists
		var exercise models.Exercise
		if err := s.db.First(&exercise, exerciseReq.ExerciseID).Error; err != nil {
			return nil, fmt.Errorf("exercise with ID %d not found", exerciseReq.ExerciseID)
		}

		workoutExercise := models.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exerciseReq.ExerciseID,
			Order:      exerciseReq.Order,
			Sets:       exerciseReq.Sets,
			Reps:       exerciseReq.Reps,
			Weight:     exerciseReq.Weight,
			Duration:   exerciseReq.Duration,
			RestTime:   exerciseReq.RestTime,
			Notes:      exerciseReq.Notes,
		}

		if err := s.db.Create(&workoutExercise).Error; err != nil {
			return nil, err
		}
	}

	// Fetch the created workout with exercises
	return s.GetWorkoutByID(workout.ID)
}

// GetWorkoutByID gets a workout by ID
func (s *WorkoutService) GetWorkoutByID(workoutID uint) (*WorkoutResponse, error) {
	var workout models.Workout
	if err := s.db.Preload("Exercises.Exercise").Preload("CreatedByUser").First(&workout, workoutID).Error; err != nil {
		return nil, errors.New("workout not found")
	}

	return s.buildWorkoutResponse(&workout), nil
}

// GetWorkouts gets all workouts with optional filtering
func (s *WorkoutService) GetWorkouts(category, difficulty, search string, createdBy *uint) ([]WorkoutResponse, error) {
	var workouts []models.Workout
	query := s.db.Preload("Exercises.Exercise").Preload("CreatedByUser")

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}
	if createdBy != nil {
		query = query.Where("created_by = ?", *createdBy)
	}

	if err := query.Where("is_active = ?", true).Order("created_at DESC").Find(&workouts).Error; err != nil {
		return nil, err
	}

	var responses []WorkoutResponse
	for _, workout := range workouts {
		responses = append(responses, *s.buildWorkoutResponse(&workout))
	}

	return responses, nil
}

// UpdateWorkout updates an existing workout
func (s *WorkoutService) UpdateWorkout(ctx context.Context, workoutID uint, req *CreateWorkoutRequest, updatedBy uint) (*WorkoutResponse, error) {
	var workout models.Workout
	if err := s.db.First(&workout, workoutID).Error; err != nil {
		return nil, errors.New("workout not found")
	}

	// Check if user has permission to update (creator or admin)
	if workout.CreatedBy != updatedBy {
		// Check if user is admin
		var user models.User
		if err := s.db.First(&user, updatedBy).Error; err != nil {
			return nil, errors.New("user not found")
		}
		if user.Role != "admin" {
			return nil, errors.New("only the creator or admin can update this workout")
		}
	}

	// Update workout fields
	workout.Name = req.Name
	workout.Description = req.Description
	workout.Category = req.Category
	workout.Difficulty = req.Difficulty
	workout.Duration = req.Duration
	workout.Equipment = req.Equipment
	workout.IsTemplate = req.IsTemplate

	if err := s.db.Save(&workout).Error; err != nil {
		return nil, err
	}

	// Update exercises (delete existing and create new ones)
	if err := s.db.Where("workout_id = ?", workoutID).Delete(&models.WorkoutExercise{}).Error; err != nil {
		return nil, err
	}

	// Add new exercises
	for _, exerciseReq := range req.Exercises {
		workoutExercise := models.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exerciseReq.ExerciseID,
			Order:      exerciseReq.Order,
			Sets:       exerciseReq.Sets,
			Reps:       exerciseReq.Reps,
			Weight:     exerciseReq.Weight,
			Duration:   exerciseReq.Duration,
			RestTime:   exerciseReq.RestTime,
			Notes:      exerciseReq.Notes,
		}

		if err := s.db.Create(&workoutExercise).Error; err != nil {
			return nil, err
		}
	}

	return s.GetWorkoutByID(workout.ID)
}

// DeleteWorkout deletes a workout (soft delete)
func (s *WorkoutService) DeleteWorkout(workoutID uint, deletedBy uint) error {
	var workout models.Workout
	if err := s.db.First(&workout, workoutID).Error; err != nil {
		return errors.New("workout not found")
	}

	// Check if user has permission to delete
	if workout.CreatedBy != deletedBy {
		// Check if user is admin
		var user models.User
		if err := s.db.First(&user, deletedBy).Error; err != nil {
			return errors.New("user not found")
		}
		if user.Role != "admin" {
			return errors.New("only the creator or admin can delete this workout")
		}
	}

	// Soft delete
	workout.IsActive = false
	return s.db.Save(&workout).Error
}

// GetWorkoutStats gets workout statistics
func (s *WorkoutService) GetWorkoutStats() (*WorkoutStats, error) {
	var stats WorkoutStats

	// Total workouts
	if err := s.db.Model(&models.Workout{}).Where("is_active = ?", true).Count(&stats.TotalWorkouts).Error; err != nil {
		return nil, err
	}

	// Active workouts (templates)
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND is_template = ?", true, true).Count(&stats.ActiveWorkouts).Error; err != nil {
		return nil, err
	}

	// Category counts
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND category = ?", true, "strength").Count(&stats.StrengthWorkouts).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND category = ?", true, "cardio").Count(&stats.CardioWorkouts).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND category = ?", true, "hiit").Count(&stats.HIITWorkouts).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND category = ?", true, "flexibility").Count(&stats.FlexibilityWorkouts).Error; err != nil {
		return nil, err
	}

	// Average duration
	if err := s.db.Model(&models.Workout{}).Where("is_active = ? AND duration > 0", true).Select("AVG(duration)").Scan(&stats.AvgDuration).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetWorkoutCategories gets all workout categories
func (s *WorkoutService) GetWorkoutCategories() ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.Workout{}).Where("is_active = ?", true).Distinct("category").Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetWorkoutDifficulties gets all workout difficulties
func (s *WorkoutService) GetWorkoutDifficulties() ([]string, error) {
	var difficulties []string
	if err := s.db.Model(&models.Workout{}).Where("is_active = ?", true).Distinct("difficulty").Pluck("difficulty", &difficulties).Error; err != nil {
		return nil, err
	}
	return difficulties, nil
}

// buildWorkoutResponse builds a workout response from a workout model
func (s *WorkoutService) buildWorkoutResponse(workout *models.Workout) *WorkoutResponse {
	response := &WorkoutResponse{
		ID:          workout.ID,
		Name:        workout.Name,
		Description: workout.Description,
		Category:    workout.Category,
		Difficulty:  workout.Difficulty,
		Duration:    workout.Duration,
		Equipment:   workout.Equipment,
		IsTemplate:  workout.IsTemplate,
		IsActive:    workout.IsActive,
		CreatedBy:   workout.CreatedBy,
		CreatedAt:   workout.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   workout.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add exercises
	if len(workout.Exercises) > 0 {
		response.Exercises = make([]WorkoutExerciseResponse, len(workout.Exercises))
		for i, exercise := range workout.Exercises {
			response.Exercises[i] = WorkoutExerciseResponse{
				ID:         exercise.ID,
				ExerciseID: exercise.ExerciseID,
				Order:      exercise.Order,
				Sets:       exercise.Sets,
				Reps:       exercise.Reps,
				Weight:     exercise.Weight,
				Duration:   exercise.Duration,
				RestTime:   exercise.RestTime,
				Notes:      exercise.Notes,
			}

			if exercise.Exercise != nil {
				response.Exercises[i].Exercise = &ExerciseInfo{
					ID:          exercise.Exercise.ID,
					Name:        exercise.Exercise.Name,
					Description: exercise.Exercise.Description,
					Category:    exercise.Exercise.Category,
					MuscleGroup: exercise.Exercise.MuscleGroup,
					Difficulty:  exercise.Exercise.Difficulty,
					Equipment:   exercise.Exercise.Equipment,
					ImageURL:    exercise.Exercise.ImageURL,
					VideoURL:    exercise.Exercise.VideoURL,
				}
			}
		}
	}

	// Add created by user info
	if workout.CreatedByUser != nil {
		response.CreatedByUser = &UserInfo{
			ID:        workout.CreatedByUser.ID,
			FirstName: workout.CreatedByUser.FirstName,
			LastName:  workout.CreatedByUser.LastName,
			Email:     workout.CreatedByUser.Email,
			Role:      workout.CreatedByUser.Role,
		}
	}

	return response
}
