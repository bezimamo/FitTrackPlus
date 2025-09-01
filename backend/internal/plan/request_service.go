package plan

import (
	"errors"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// PlanRequestService handles plan request management business logic
type PlanRequestService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewPlanRequestService creates a new plan request service
func NewPlanRequestService(cfg *config.Config) *PlanRequestService {
	return &PlanRequestService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// Request DTOs
type CreatePlanRequestRequest struct {
	PlanID uint   `json:"plan_id" binding:"required"`
	Reason string `json:"reason,omitempty"`
}

type PlanRequestResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	PlanID      uint      `json:"plan_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  *uint     `json:"reviewed_by,omitempty"`
	User        *UserInfo `json:"user,omitempty"`
	Plan        *PlanResponse `json:"plan,omitempty"`
	ReviewedByUser *UserInfo `json:"reviewed_by_user,omitempty"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type ApprovalRequest struct {
	TrainerID uint   `json:"trainer_id" binding:"required"`
	Notes     string `json:"notes,omitempty"`
}

// CreatePlanRequest creates a new plan request (member requests a plan)
func (s *PlanRequestService) CreatePlanRequest(userID uint, req *CreatePlanRequestRequest) (*PlanRequestResponse, error) {
	// Check if plan exists and is active
	var plan models.Plan
	if err := s.db.First(&plan, req.PlanID).Error; err != nil {
		return nil, errors.New("plan not found")
	}
	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	// Check if user already has a pending request for this plan
	var existingRequest models.PlanRequest
	err := s.db.Where("user_id = ? AND plan_id = ? AND status = ?", userID, req.PlanID, "pending").First(&existingRequest).Error
	if err == nil {
		return nil, errors.New("you already have a pending request for this plan")
	}

	// Create the request
	planRequest := models.PlanRequest{
		UserID:      userID,
		PlanID:      req.PlanID,
		Reason:      req.Reason,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.db.Create(&planRequest).Error; err != nil {
		return nil, err
	}

	// Load relationships for response
	if err := s.db.Preload("User").Preload("Plan").First(&planRequest, planRequest.ID).Error; err != nil {
		return nil, err
	}

	return s.buildPlanRequestResponse(&planRequest), nil
}

// GetPendingRequests gets all pending plan requests (admin only)
func (s *PlanRequestService) GetPendingRequests() ([]PlanRequestResponse, error) {
	var requests []models.PlanRequest
	if err := s.db.Where("status = ?", "pending").
		Preload("User").
		Preload("Plan").
		Order("requested_at DESC").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	var responses []PlanRequestResponse
	for _, request := range requests {
		responses = append(responses, *s.buildPlanRequestResponse(&request))
	}

	return responses, nil
}

// GetAllRequests gets all plan requests with optional status filter (admin only)
func (s *PlanRequestService) GetAllRequests(status string) ([]PlanRequestResponse, error) {
	var requests []models.PlanRequest
	query := s.db.Preload("User").Preload("Plan").Preload("ReviewedByUser")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("requested_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}

	var responses []PlanRequestResponse
	for _, request := range requests {
		responses = append(responses, *s.buildPlanRequestResponse(&request))
	}

	return responses, nil
}

// ApproveRequest approves a plan request and assigns trainer (admin only)
func (s *PlanRequestService) ApproveRequest(requestID uint, adminID uint, approval *ApprovalRequest) (*PlanRequestResponse, error) {
	// Get the request
	var planRequest models.PlanRequest
	if err := s.db.Preload("User").Preload("Plan").First(&planRequest, requestID).Error; err != nil {
		return nil, errors.New("request not found")
	}

	if planRequest.Status != "pending" {
		return nil, errors.New("request is not pending")
	}

	// Verify trainer exists and has trainer role
	var trainer models.User
	if err := s.db.First(&trainer, approval.TrainerID).Error; err != nil {
		return nil, errors.New("trainer not found")
	}
	if trainer.Role != "trainer" {
		return nil, errors.New("assigned user is not a trainer")
	}

	// Start transaction
	tx := s.db.Begin()

	// Update request status
	now := time.Now()
	planRequest.Status = "approved"
	planRequest.ReviewedAt = &now
	planRequest.ReviewedBy = &adminID

	if err := tx.Save(&planRequest).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create UserPlan (assign plan to member)
	userPlan := models.UserPlan{
		UserID:     planRequest.UserID,
		PlanID:     planRequest.PlanID,
		Status:     "active",
		AssignedAt: time.Now(),
	}

	if err := tx.Create(&userPlan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create TrainerAssignment (assign trainer to manage the plan)
	trainerAssignment := models.TrainerAssignment{
		TrainerID:  approval.TrainerID,
		MemberID:   planRequest.UserID,
		PlanID:     planRequest.PlanID,
		UserPlanID: userPlan.ID,
		Status:     "active",
		AssignedAt: time.Now(),
		AssignedBy: adminID,
		Notes:      approval.Notes,
	}

	if err := tx.Create(&trainerAssignment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("User").Preload("Plan").Preload("ReviewedByUser").First(&planRequest, planRequest.ID).Error; err != nil {
		return nil, err
	}

	return s.buildPlanRequestResponse(&planRequest), nil
}

// RejectRequest rejects a plan request (admin only)
func (s *PlanRequestService) RejectRequest(requestID uint, adminID uint) (*PlanRequestResponse, error) {
	var planRequest models.PlanRequest
	if err := s.db.First(&planRequest, requestID).Error; err != nil {
		return nil, errors.New("request not found")
	}

	if planRequest.Status != "pending" {
		return nil, errors.New("request is not pending")
	}

	// Update request status
	now := time.Now()
	planRequest.Status = "rejected"
	planRequest.ReviewedAt = &now
	planRequest.ReviewedBy = &adminID

	if err := s.db.Save(&planRequest).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("User").Preload("Plan").Preload("ReviewedByUser").First(&planRequest, planRequest.ID).Error; err != nil {
		return nil, err
	}

	return s.buildPlanRequestResponse(&planRequest), nil
}

// GetUserRequests gets all requests made by a specific user
func (s *PlanRequestService) GetUserRequests(userID uint) ([]PlanRequestResponse, error) {
	var requests []models.PlanRequest
	if err := s.db.Where("user_id = ?", userID).
		Preload("Plan").
		Preload("ReviewedByUser").
		Order("requested_at DESC").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	var responses []PlanRequestResponse
	for _, request := range requests {
		responses = append(responses, *s.buildPlanRequestResponse(&request))
	}

	return responses, nil
}

// Helper function to build response
func (s *PlanRequestService) buildPlanRequestResponse(request *models.PlanRequest) *PlanRequestResponse {
	response := &PlanRequestResponse{
		ID:          request.ID,
		UserID:      request.UserID,
		PlanID:      request.PlanID,
		Reason:      request.Reason,
		Status:      request.Status,
		RequestedAt: request.RequestedAt,
		ReviewedAt:  request.ReviewedAt,
		ReviewedBy:  request.ReviewedBy,
	}

	// Add user info
	if request.User.ID != 0 {
		response.User = &UserInfo{
			ID:        request.User.ID,
			FirstName: request.User.FirstName,
			LastName:  request.User.LastName,
			Email:     request.User.Email,
			Role:      request.User.Role,
		}
	}

	// Add plan info
	if request.Plan.ID != 0 {
		response.Plan = &PlanResponse{
			ID:          request.Plan.ID,
			Name:        request.Plan.Name,
			Description: request.Plan.Description,
			GoalType:    request.Plan.GoalType,
			PlanType:    request.Plan.PlanType,
			Duration:    request.Plan.Duration,
			IsActive:    request.Plan.IsActive,
			CreatedAt:   request.Plan.CreatedAt,
			UpdatedAt:   request.Plan.UpdatedAt,
		}
	}

	// Add reviewed by info
	if request.ReviewedByUser != nil && request.ReviewedByUser.ID != 0 {
		response.ReviewedByUser = &UserInfo{
			ID:        request.ReviewedByUser.ID,
			FirstName: request.ReviewedByUser.FirstName,
			LastName:  request.ReviewedByUser.LastName,
			Email:     request.ReviewedByUser.Email,
			Role:      request.ReviewedByUser.Role,
		}
	}

	return response
}
