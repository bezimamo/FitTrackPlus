package plan

import (
	"errors"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// PlanService handles plan management business logic
type PlanService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewPlanService creates a new plan service
func NewPlanService(cfg *config.Config) *PlanService {
	return &PlanService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// PlanRequest represents a plan creation/update request
type PlanRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	GoalType    string `json:"goal_type" binding:"required"` // lose_weight, gain_muscle, flexibility, rehab
	PlanType    string `json:"plan_type" binding:"required"` // fitness, diet, physio
	Duration    int    `json:"duration" binding:"min=1"`     // in days, minimum 1 day
}

// PlanResponse represents a plan response
type PlanResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GoalType    string    `json:"goal_type"`
	PlanType    string    `json:"plan_type"`
	Duration    int       `json:"duration"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserPlanResponse represents a user's assigned plan
type UserPlanResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	PlanID       uint      `json:"plan_id"`
	Status       string    `json:"status"` // active, completed, paused
	AssignedAt   time.Time `json:"assigned_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Progress     float64   `json:"progress"` // 0-100
	Plan         PlanResponse `json:"plan"`
}

// CreatePlan creates a new plan template
func (s *PlanService) CreatePlan(req *PlanRequest, createdBy uint) (*PlanResponse, error) {
	// Validate goal type
	if !isValidGoalType(req.GoalType) {
		return nil, errors.New("invalid goal type: must be one of lose_weight, gain_muscle, flexibility, rehab")
	}

	// Validate plan type
	if !isValidPlanType(req.PlanType) {
		return nil, errors.New("invalid plan type: must be one of fitness, diet, physio")
	}

	// Validate duration
	if req.Duration < 1 {
		return nil, errors.New("duration must be at least 1 day")
	}

	// Create plan
	plan := models.Plan{
		Name:        req.Name,
		Description: req.Description,
		GoalType:    req.GoalType,
		PlanType:    req.PlanType,
		Duration:    req.Duration,
		IsActive:    true,
	}

	if err := s.db.Create(&plan).Error; err != nil {
		return nil, err
	}

	return s.buildPlanResponse(&plan), nil
}

// GetPlans retrieves all plans with optional filtering
func (s *PlanService) GetPlans(goalType, planType string, isActive *bool) ([]PlanResponse, error) {
	var plans []models.Plan
	query := s.db

	// Apply filters
	if goalType != "" {
		query = query.Where("goal_type = ?", goalType)
	}
	if planType != "" {
		query = query.Where("plan_type = ?", planType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []PlanResponse
	for _, plan := range plans {
		responses = append(responses, *s.buildPlanResponse(&plan))
	}

	return responses, nil
}

// GetPlan retrieves a specific plan by ID
func (s *PlanService) GetPlan(planID uint) (*PlanResponse, error) {
	var plan models.Plan
	if err := s.db.First(&plan, planID).Error; err != nil {
		return nil, err
	}

	return s.buildPlanResponse(&plan), nil
}

// AssignPlan assigns a plan to a user
func (s *PlanService) AssignPlan(userID, planID uint, assignedBy uint) (*UserPlanResponse, error) {
	// Check if plan exists and is active
	var plan models.Plan
	if err := s.db.First(&plan, planID).Error; err != nil {
		return nil, errors.New("plan not found")
	}
	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	// Check if user already has an active plan
	var existingUserPlan models.UserPlan
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&existingUserPlan).Error
	if err == nil {
		// User has an active plan, pause it
		existingUserPlan.Status = "paused"
		s.db.Save(&existingUserPlan)
	}

	// Create new user plan assignment
	userPlan := models.UserPlan{
		UserID:     userID,
		PlanID:     planID,
		Status:     "active",
		AssignedAt: time.Now(),
	}

	if err := s.db.Create(&userPlan).Error; err != nil {
		return nil, err
	}

	return s.buildUserPlanResponse(&userPlan), nil
}

// GetUserPlans retrieves plans assigned to a specific user
func (s *PlanService) GetUserPlans(userID uint) ([]UserPlanResponse, error) {
	var userPlans []models.UserPlan
	if err := s.db.Where("user_id = ?", userID).Preload("Plan").Find(&userPlans).Error; err != nil {
		return nil, err
	}

	var responses []UserPlanResponse
	for _, userPlan := range userPlans {
		response := s.buildUserPlanResponse(&userPlan)
		responses = append(responses, *response)
	}

	return responses, nil
}

// Helper methods
func (s *PlanService) buildPlanResponse(plan *models.Plan) *PlanResponse {
	return &PlanResponse{
		ID:          plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
		GoalType:    plan.GoalType,
		PlanType:    plan.PlanType,
		Duration:    plan.Duration,
		IsActive:    plan.IsActive,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}

func (s *PlanService) buildUserPlanResponse(userPlan *models.UserPlan) *UserPlanResponse {
	response := &UserPlanResponse{
		ID:          userPlan.ID,
		UserID:      userPlan.UserID,
		PlanID:      userPlan.PlanID,
		Status:      userPlan.Status,
		AssignedAt:  userPlan.AssignedAt,
		CompletedAt: userPlan.CompletedAt,
		Progress:    s.calculatePlanProgress(userPlan),
	}

	// Load plan details
	if userPlan.Plan.ID != 0 {
		response.Plan = *s.buildPlanResponse(&userPlan.Plan)
	}

	return response
}

func (s *PlanService) calculatePlanProgress(userPlan *models.UserPlan) float64 {
	// For now, return a simple calculation
	// In a real app, you'd calculate based on completed workouts/sessions
	if userPlan.Status == "completed" {
		return 100.0
	} else if userPlan.Status == "paused" {
		return 50.0 // Example: paused plans show 50% progress
	}
	return 0.0 // Active plans start at 0%
}

// Validation helpers
func isValidGoalType(goalType string) bool {
	validTypes := []string{"lose_weight", "gain_muscle", "flexibility", "rehab", "weight_loss"}
	for _, t := range validTypes {
		if t == goalType {
			return true
		}
	}
	return false
}

func isValidPlanType(planType string) bool {
	validTypes := []string{"fitness", "diet", "physio", "weight"}
	for _, t := range validTypes {
		if t == planType {
			return true
		}
	}
	return false
}
