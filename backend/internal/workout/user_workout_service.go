package workout

import (
	"context"
	"errors"
	"time"
	"fittrackplus/internal/common/models"
	"gorm.io/gorm"
)

// UserWorkoutService handles user workout session tracking
type UserWorkoutService struct {
	db *gorm.DB
}

// NewUserWorkoutService creates a new user workout service
func NewUserWorkoutService(db *gorm.DB) *UserWorkoutService {
	return &UserWorkoutService{db: db}
}

// StartWorkoutRequest represents the request to start a workout
type StartWorkoutRequest struct {
	WorkoutID uint   `json:"workout_id" binding:"required"`
	Notes     string `json:"notes"`
}

// UserWorkoutResponse represents a user workout session
type UserWorkoutResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	WorkoutID   uint      `json:"workout_id"`
	StartedAt   string    `json:"started_at"`
	CompletedAt *string   `json:"completed_at,omitempty"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	Workout     *WorkoutResponse `json:"workout,omitempty"`
	User        *UserInfo `json:"user,omitempty"`
	ExerciseSets []UserWorkoutExerciseResponse `json:"exercise_sets,omitempty"`
}

// UserWorkoutExerciseResponse represents a user's exercise performance
type UserWorkoutExerciseResponse struct {
	ID            uint      `json:"id"`
	UserWorkoutID uint      `json:"user_workout_id"`
	ExerciseID    uint      `json:"exercise_id"`
	SetNumber     int       `json:"set_number"`
	Reps          int       `json:"reps"`
	Weight        float64   `json:"weight"`
	Duration      int       `json:"duration"`
	RestTime      int       `json:"rest_time"`
	Notes         string    `json:"notes"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Exercise      *ExerciseInfo `json:"exercise,omitempty"`
}

// CompleteExerciseSetRequest represents completing an exercise set
type CompleteExerciseSetRequest struct {
	ExerciseID uint    `json:"exercise_id" binding:"required"`
	SetNumber  int     `json:"set_number" binding:"required"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Duration   int     `json:"duration"`
	RestTime   int     `json:"rest_time"`
	Notes      string  `json:"notes"`
}

// WorkoutProgressResponse represents workout progress information
type WorkoutProgressResponse struct {
	UserWorkoutID     uint    `json:"user_workout_id"`
	WorkoutName       string  `json:"workout_name"`
	TotalExercises    int     `json:"total_exercises"`
	CompletedExercises int    `json:"completed_exercises"`
	TotalSets         int     `json:"total_sets"`
	CompletedSets     int     `json:"completed_sets"`
	ProgressPercent   float64 `json:"progress_percent"`
	ElapsedTime       int     `json:"elapsed_time"` // in minutes
	Status            string  `json:"status"`
}

// StartWorkout starts a new workout session for a user
func (s *UserWorkoutService) StartWorkout(ctx context.Context, req *StartWorkoutRequest, userID uint) (*UserWorkoutResponse, error) {
	// Verify workout exists and is active
	var workout models.Workout
	if err := s.db.Where("id = ? AND is_active = ?", req.WorkoutID, true).First(&workout).Error; err != nil {
		return nil, errors.New("workout not found or inactive")
	}

	// Check if user already has an active workout
	var existingWorkout models.UserWorkout
	if err := s.db.Where("user_id = ? AND status = ?", userID, "in_progress").First(&existingWorkout).Error; err == nil {
		return nil, errors.New("user already has an active workout session")
	}

	// Create new workout session
	userWorkout := models.UserWorkout{
		UserID:    userID,
		WorkoutID: req.WorkoutID,
		StartedAt: time.Now(),
		Status:    "in_progress",
		Notes:     req.Notes,
	}

	if err := s.db.Create(&userWorkout).Error; err != nil {
		return nil, err
	}

	return s.GetUserWorkoutByID(userWorkout.ID)
}

// CompleteExerciseSet records a completed exercise set
func (s *UserWorkoutService) CompleteExerciseSet(ctx context.Context, userWorkoutID uint, req *CompleteExerciseSetRequest, userID uint) (*UserWorkoutExerciseResponse, error) {
	// Verify user workout exists and belongs to user
	var userWorkout models.UserWorkout
	if err := s.db.Where("id = ? AND user_id = ?", userWorkoutID, userID).First(&userWorkout).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	// Verify workout is still in progress
	if userWorkout.Status != "in_progress" {
		return nil, errors.New("workout session is not in progress")
	}

	// Verify exercise exists in the workout
	var workoutExercise models.WorkoutExercise
	if err := s.db.Where("workout_id = ? AND exercise_id = ?", userWorkout.WorkoutID, req.ExerciseID).First(&workoutExercise).Error; err != nil {
		return nil, errors.New("exercise not found in this workout")
	}

	// Create exercise set record
	exerciseSet := models.UserWorkoutExercise{
		UserWorkoutID: userWorkoutID,
		ExerciseID:    req.ExerciseID,
		SetNumber:     req.SetNumber,
		Reps:          req.Reps,
		Weight:        req.Weight,
		Duration:      req.Duration,
		RestTime:      req.RestTime,
		Notes:         req.Notes,
	}

	if err := s.db.Create(&exerciseSet).Error; err != nil {
		return nil, err
	}

	return s.buildUserWorkoutExerciseResponse(&exerciseSet), nil
}

// CompleteWorkout marks a workout as completed
func (s *UserWorkoutService) CompleteWorkout(ctx context.Context, userWorkoutID uint, userID uint) (*UserWorkoutResponse, error) {
	// Verify user workout exists and belongs to user
	var userWorkout models.UserWorkout
	if err := s.db.Where("id = ? AND user_id = ?", userWorkoutID, userID).First(&userWorkout).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	// Verify workout is in progress
	if userWorkout.Status != "in_progress" {
		return nil, errors.New("workout session is not in progress")
	}

	// Update workout status
	now := time.Now()
	userWorkout.Status = "completed"
	userWorkout.CompletedAt = &now

	if err := s.db.Save(&userWorkout).Error; err != nil {
		return nil, err
	}

	return s.GetUserWorkoutByID(userWorkoutID)
}

// PauseWorkout pauses a workout session
func (s *UserWorkoutService) PauseWorkout(ctx context.Context, userWorkoutID uint, userID uint) (*UserWorkoutResponse, error) {
	// Verify user workout exists and belongs to user
	var userWorkout models.UserWorkout
	if err := s.db.Where("id = ? AND user_id = ?", userWorkoutID, userID).First(&userWorkout).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	// Verify workout is in progress
	if userWorkout.Status != "in_progress" {
		return nil, errors.New("workout session is not in progress")
	}

	// Update workout status
	userWorkout.Status = "paused"

	if err := s.db.Save(&userWorkout).Error; err != nil {
		return nil, err
	}

	return s.GetUserWorkoutByID(userWorkoutID)
}

// ResumeWorkout resumes a paused workout session
func (s *UserWorkoutService) ResumeWorkout(ctx context.Context, userWorkoutID uint, userID uint) (*UserWorkoutResponse, error) {
	// Verify user workout exists and belongs to user
	var userWorkout models.UserWorkout
	if err := s.db.Where("id = ? AND user_id = ?", userWorkoutID, userID).First(&userWorkout).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	// Verify workout is paused
	if userWorkout.Status != "paused" {
		return nil, errors.New("workout session is not paused")
	}

	// Update workout status
	userWorkout.Status = "in_progress"

	if err := s.db.Save(&userWorkout).Error; err != nil {
		return nil, err
	}

	return s.GetUserWorkoutByID(userWorkoutID)
}

// GetUserWorkoutByID gets a user workout by ID
func (s *UserWorkoutService) GetUserWorkoutByID(userWorkoutID uint) (*UserWorkoutResponse, error) {
	var userWorkout models.UserWorkout
	if err := s.db.Preload("Workout.Exercises.Exercise").Preload("User").Preload("ExerciseSets.Exercise").First(&userWorkout, userWorkoutID).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	return s.buildUserWorkoutResponse(&userWorkout), nil
}

// GetUserWorkouts gets all workouts for a specific user
func (s *UserWorkoutService) GetUserWorkouts(userID uint, status string) ([]UserWorkoutResponse, error) {
	var userWorkouts []models.UserWorkout
	query := s.db.Preload("Workout").Preload("User").Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&userWorkouts).Error; err != nil {
		return nil, err
	}

	var responses []UserWorkoutResponse
	for _, userWorkout := range userWorkouts {
		responses = append(responses, *s.buildUserWorkoutResponse(&userWorkout))
	}

	return responses, nil
}

// GetWorkoutProgress gets the progress of a specific workout
func (s *UserWorkoutService) GetWorkoutProgress(userWorkoutID uint, userID uint) (*WorkoutProgressResponse, error) {
	// Get user workout
	var userWorkout models.UserWorkout
	if err := s.db.Preload("Workout.Exercises").Where("id = ? AND user_id = ?", userWorkoutID, userID).First(&userWorkout).Error; err != nil {
		return nil, errors.New("workout session not found")
	}

	// Get completed exercise sets
	var completedSets int64
	if err := s.db.Model(&models.UserWorkoutExercise{}).Where("user_workout_id = ?", userWorkoutID).Count(&completedSets).Error; err != nil {
		return nil, err
	}

	// Calculate total sets from workout exercises
	totalSets := 0
	for _, exercise := range userWorkout.Workout.Exercises {
		totalSets += exercise.Sets
	}

	// Calculate progress
	progressPercent := 0.0
	if totalSets > 0 {
		progressPercent = float64(completedSets) / float64(totalSets) * 100
	}

	// Calculate elapsed time
	elapsedTime := int(time.Since(userWorkout.StartedAt).Minutes())

	return &WorkoutProgressResponse{
		UserWorkoutID:     userWorkout.ID,
		WorkoutName:       userWorkout.Workout.Name,
		TotalExercises:    len(userWorkout.Workout.Exercises),
		CompletedExercises: 0, // TODO: Calculate based on completed sets
		TotalSets:         totalSets,
		CompletedSets:     int(completedSets),
		ProgressPercent:   progressPercent,
		ElapsedTime:       elapsedTime,
		Status:            userWorkout.Status,
	}, nil
}

// buildUserWorkoutResponse builds a user workout response from a model
func (s *UserWorkoutService) buildUserWorkoutResponse(userWorkout *models.UserWorkout) *UserWorkoutResponse {
	response := &UserWorkoutResponse{
		ID:        userWorkout.ID,
		UserID:    userWorkout.UserID,
		WorkoutID: userWorkout.WorkoutID,
		StartedAt: userWorkout.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		Status:    userWorkout.Status,
		Notes:     userWorkout.Notes,
		CreatedAt: userWorkout.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: userWorkout.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if userWorkout.CompletedAt != nil {
		completedAt := userWorkout.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	// Add workout info
	if userWorkout.Workout != nil {
		response.Workout = &WorkoutResponse{
			ID:          userWorkout.Workout.ID,
			Name:        userWorkout.Workout.Name,
			Description: userWorkout.Workout.Description,
			Category:    userWorkout.Workout.Category,
			Difficulty:  userWorkout.Workout.Difficulty,
			Duration:    userWorkout.Workout.Duration,
			Equipment:   userWorkout.Workout.Equipment,
			IsTemplate:  userWorkout.Workout.IsTemplate,
			IsActive:    userWorkout.Workout.IsActive,
			CreatedBy:   userWorkout.Workout.CreatedBy,
			CreatedAt:   userWorkout.Workout.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   userWorkout.Workout.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Add user info
	if userWorkout.User != nil {
		response.User = &UserInfo{
			ID:        userWorkout.User.ID,
			FirstName: userWorkout.User.FirstName,
			LastName:  userWorkout.User.LastName,
			Email:     userWorkout.User.Email,
			Role:      userWorkout.User.Role,
		}
	}

	// Add exercise sets
	if len(userWorkout.ExerciseSets) > 0 {
		response.ExerciseSets = make([]UserWorkoutExerciseResponse, len(userWorkout.ExerciseSets))
		for i, exerciseSet := range userWorkout.ExerciseSets {
			response.ExerciseSets[i] = *s.buildUserWorkoutExerciseResponse(&exerciseSet)
		}
	}

	return response
}

// buildUserWorkoutExerciseResponse builds a user workout exercise response from a model
func (s *UserWorkoutService) buildUserWorkoutExerciseResponse(exerciseSet *models.UserWorkoutExercise) *UserWorkoutExerciseResponse {
	response := &UserWorkoutExerciseResponse{
		ID:            exerciseSet.ID,
		UserWorkoutID: exerciseSet.UserWorkoutID,
		ExerciseID:    exerciseSet.ExerciseID,
		SetNumber:     exerciseSet.SetNumber,
		Reps:          exerciseSet.Reps,
		Weight:        exerciseSet.Weight,
		Duration:      exerciseSet.Duration,
		RestTime:      exerciseSet.RestTime,
		Notes:         exerciseSet.Notes,
		CreatedAt:     exerciseSet.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     exerciseSet.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add exercise info
	if exerciseSet.Exercise != nil {
		response.Exercise = &ExerciseInfo{
			ID:          exerciseSet.Exercise.ID,
			Name:        exerciseSet.Exercise.Name,
			Description: exerciseSet.Exercise.Description,
			Category:    exerciseSet.Exercise.Category,
			MuscleGroup: exerciseSet.Exercise.MuscleGroup,
			Difficulty:  exerciseSet.Exercise.Difficulty,
			Equipment:   exerciseSet.Exercise.Equipment,
			ImageURL:    exerciseSet.Exercise.ImageURL,
			VideoURL:    exerciseSet.Exercise.VideoURL,
		}
	}

	return response
}
