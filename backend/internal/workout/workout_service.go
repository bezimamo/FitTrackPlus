package workout

import (
	"errors"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// WorkoutService handles workout management business logic
type WorkoutService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(cfg *config.Config) *WorkoutService {
	return &WorkoutService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// Request DTOs
type CreateWorkoutRequest struct {
	Name              string                    `json:"name" binding:"required"`
	Description       string                    `json:"description"`
	MemberID          uint                      `json:"member_id" binding:"required"`
	PlanID            *uint                     `json:"plan_id,omitempty"`
	Difficulty        string                    `json:"difficulty" binding:"required"`
	EstimatedDuration int                       `json:"estimated_duration" binding:"required,min=1"`
	Notes             string                    `json:"notes"`
	Exercises         []CreateExerciseRequest   `json:"exercises" binding:"required,min=1"`
}

type CreateExerciseRequest struct {
	ExerciseID uint     `json:"exercise_id" binding:"required"`
	Order      int      `json:"order" binding:"required,min=1"`
	Sets       int      `json:"sets" binding:"required,min=1"`
	Reps       int      `json:"reps" binding:"required,min=1"`
	Weight     *float64 `json:"weight,omitempty"`
	Duration   *int     `json:"duration,omitempty"`
	RestTime   int      `json:"rest_time" binding:"required,min=0"`
	Notes      string   `json:"notes"`
}

type WorkoutResponse struct {
	ID                uint                     `json:"id"`
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	MemberID          uint                     `json:"member_id"`
	TrainerID         uint                     `json:"trainer_id"`
	PlanID            *uint                    `json:"plan_id,omitempty"`
	Difficulty        string                   `json:"difficulty"`
	EstimatedDuration int                      `json:"estimated_duration"`
	Notes             string                   `json:"notes"`
	IsActive          bool                     `json:"is_active"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	
	// Relationships
	Member            *UserInfo                `json:"member,omitempty"`
	Trainer           *UserInfo                `json:"trainer,omitempty"`
	Plan              *PlanInfo                `json:"plan,omitempty"`
	Exercises         []WorkoutExerciseResponse `json:"exercises,omitempty"`
}

type WorkoutExerciseResponse struct {
	ID         uint     `json:"id"`
	Order      int      `json:"order"`
	Sets       int      `json:"sets"`
	Reps       int      `json:"reps"`
	Weight     *float64 `json:"weight,omitempty"`
	Duration   *int     `json:"duration,omitempty"`
	RestTime   int      `json:"rest_time"`
	Notes      string   `json:"notes"`
	
	// Exercise details
	Exercise   ExerciseInfo `json:"exercise"`
}

type ExerciseInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MuscleGroup string `json:"muscle_group"`
	Difficulty  string `json:"difficulty"`
	Equipment   string `json:"equipment"`
	VideoURL    string `json:"video_url"`
	ImageURL    string `json:"image_url"`
	Instructions string `json:"instructions"`
	Tips        string `json:"tips"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type PlanInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GoalType    string `json:"goal_type"`
	PlanType    string `json:"plan_type"`
}

// CreateWorkout creates a new workout for a member (trainer only)
func (s *WorkoutService) CreateWorkout(trainerID uint, req *CreateWorkoutRequest) (*WorkoutResponse, error) {
	// Verify trainer has assignment to this member
	var assignment models.TrainerAssignment
	if err := s.db.Where("trainer_id = ? AND member_id = ? AND status = ?", trainerID, req.MemberID, "active").First(&assignment).Error; err != nil {
		return nil, errors.New("you are not assigned to this member or assignment is not active")
	}

	// Start transaction
	tx := s.db.Begin()

	// Create the workout
	workout := models.Workout{
		Name:              req.Name,
		Description:       req.Description,
		MemberID:          req.MemberID,
		TrainerID:         trainerID,
		PlanID:            req.PlanID,
		Difficulty:        req.Difficulty,
		EstimatedDuration: req.EstimatedDuration,
		Notes:             req.Notes,
		IsActive:          true,
	}

	if err := tx.Create(&workout).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create workout exercises
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

		if err := tx.Create(&workoutExercise).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Member").Preload("Trainer").Preload("Plan").Preload("Exercises.Exercise").First(&workout, workout.ID).Error; err != nil {
		return nil, err
	}

	return s.buildWorkoutResponse(&workout), nil
}

// GetTrainerWorkouts gets all workouts created by a trainer
func (s *WorkoutService) GetTrainerWorkouts(trainerID uint, memberID *uint, status string) ([]WorkoutResponse, error) {
	var workouts []models.Workout
	query := s.db.Where("trainer_id = ?", trainerID).
		Preload("Member").
		Preload("Plan").
		Preload("Exercises.Exercise")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}

	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Order("created_at DESC").Find(&workouts).Error; err != nil {
		return nil, err
	}

	var responses []WorkoutResponse
	for _, workout := range workouts {
		responses = append(responses, *s.buildWorkoutResponse(&workout))
	}

	return responses, nil
}

// GetMemberWorkouts gets all workouts assigned to a member
func (s *WorkoutService) GetMemberWorkouts(memberID uint, status string) ([]WorkoutResponse, error) {
	var workouts []models.Workout
	query := s.db.Where("member_id = ?", memberID).
		Preload("Trainer").
		Preload("Plan").
		Preload("Exercises.Exercise")

	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Order("created_at DESC").Find(&workouts).Error; err != nil {
		return nil, err
	}

	var responses []WorkoutResponse
	for _, workout := range workouts {
		responses = append(responses, *s.buildWorkoutResponse(&workout))
	}

	return responses, nil
}

// GetWorkout gets a specific workout with all details
func (s *WorkoutService) GetWorkout(workoutID uint, userID uint, userRole string) (*WorkoutResponse, error) {
	var workout models.Workout
	if err := s.db.Preload("Member").Preload("Trainer").Preload("Plan").Preload("Exercises.Exercise").First(&workout, workoutID).Error; err != nil {
		return nil, err
	}

	// Check access permissions
	if userRole == "trainer" && workout.TrainerID != userID {
		return nil, errors.New("you can only view workouts you created")
	}
	if userRole == "member" && workout.MemberID != userID {
		return nil, errors.New("you can only view workouts assigned to you")
	}

	return s.buildWorkoutResponse(&workout), nil
}

// UpdateWorkout updates an existing workout (trainer only)
func (s *WorkoutService) UpdateWorkout(workoutID uint, trainerID uint, req *CreateWorkoutRequest) (*WorkoutResponse, error) {
	var workout models.Workout
	if err := s.db.First(&workout, workoutID).Error; err != nil {
		return nil, errors.New("workout not found")
	}

	// Verify trainer owns this workout
	if workout.TrainerID != trainerID {
		return nil, errors.New("you can only update workouts you created")
	}

	// Start transaction
	tx := s.db.Begin()

	// Update workout details
	workout.Name = req.Name
	workout.Description = req.Description
	workout.Difficulty = req.Difficulty
	workout.EstimatedDuration = req.EstimatedDuration
	workout.Notes = req.Notes

	if err := tx.Save(&workout).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete existing exercises
	if err := tx.Where("workout_id = ?", workoutID).Delete(&models.WorkoutExercise{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new workout exercises
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

		if err := tx.Create(&workoutExercise).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Member").Preload("Trainer").Preload("Plan").Preload("Exercises.Exercise").First(&workout, workout.ID).Error; err != nil {
		return nil, err
	}

	return s.buildWorkoutResponse(&workout), nil
}

// Helper function to build response
func (s *WorkoutService) buildWorkoutResponse(workout *models.Workout) *WorkoutResponse {
	response := &WorkoutResponse{
		ID:                workout.ID,
		Name:              workout.Name,
		Description:       workout.Description,
		MemberID:          workout.MemberID,
		TrainerID:         workout.TrainerID,
		PlanID:            workout.PlanID,
		Difficulty:        workout.Difficulty,
		EstimatedDuration: workout.EstimatedDuration,
		Notes:             workout.Notes,
		IsActive:          workout.IsActive,
		CreatedAt:         workout.CreatedAt,
		UpdatedAt:         workout.UpdatedAt,
	}

	// Add member info
	if workout.Member.ID != 0 {
		response.Member = &UserInfo{
			ID:        workout.Member.ID,
			FirstName: workout.Member.FirstName,
			LastName:  workout.Member.LastName,
			Email:     workout.Member.Email,
			Role:      workout.Member.Role,
		}
	}

	// Add trainer info
	if workout.Trainer.ID != 0 {
		response.Trainer = &UserInfo{
			ID:        workout.Trainer.ID,
			FirstName: workout.Trainer.FirstName,
			LastName:  workout.Trainer.LastName,
			Email:     workout.Trainer.Email,
			Role:      workout.Trainer.Role,
		}
	}

	// Add plan info
	if workout.Plan != nil && workout.Plan.ID != 0 {
		response.Plan = &PlanInfo{
			ID:          workout.Plan.ID,
			Name:        workout.Plan.Name,
			Description: workout.Plan.Description,
			GoalType:    workout.Plan.GoalType,
			PlanType:    workout.Plan.PlanType,
		}
	}

	// Add exercises
	if len(workout.Exercises) > 0 {
		response.Exercises = make([]WorkoutExerciseResponse, len(workout.Exercises))
		for i, exercise := range workout.Exercises {
			response.Exercises[i] = WorkoutExerciseResponse{
				ID:       exercise.ID,
				Order:    exercise.Order,
				Sets:     exercise.Sets,
				Reps:     exercise.Reps,
				Weight:   exercise.Weight,
				Duration: exercise.Duration,
				RestTime: exercise.RestTime,
				Notes:    exercise.Notes,
				Exercise: ExerciseInfo{
					ID:          exercise.Exercise.ID,
					Name:        exercise.Exercise.Name,
					Description: exercise.Exercise.Description,
					Category:    exercise.Exercise.Category,
					MuscleGroup: exercise.Exercise.MuscleGroup,
					Difficulty:  exercise.Exercise.Difficulty,
					Equipment:   exercise.Exercise.Equipment,
					VideoURL:    exercise.Exercise.VideoURL,
					ImageURL:    exercise.Exercise.ImageURL,
					Instructions: exercise.Exercise.Instructions,
					Tips:        exercise.Exercise.Tips,
				},
			}
		}
	}

	return response
}
