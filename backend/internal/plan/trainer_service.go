package plan

import (
	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// TrainerService handles trainer assignment management
type TrainerService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewTrainerService creates a new trainer service
func NewTrainerService(cfg *config.Config) *TrainerService {
	return &TrainerService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// TrainerAssignmentResponse represents a trainer assignment
type TrainerAssignmentResponse struct {
	ID         uint      `json:"id"`
	TrainerID  uint      `json:"trainer_id"`
	MemberID   uint      `json:"member_id"`
	PlanID     uint      `json:"plan_id"`
	UserPlanID uint      `json:"user_plan_id"`
	Status     string    `json:"status"`
	AssignedAt string    `json:"assigned_at"`
	AssignedBy uint      `json:"assigned_by"`
	Notes      string    `json:"notes"`
	
	// Detailed information
	Trainer      *UserInfo        `json:"trainer,omitempty"`
	Member       *UserInfo        `json:"member,omitempty"`
	Plan         *PlanResponse    `json:"plan,omitempty"`
	UserPlan     *UserPlanResponse `json:"user_plan,omitempty"`
	AssignedByUser *UserInfo      `json:"assigned_by_user,omitempty"`
}

// GetTrainerAssignments gets all assignments for a specific trainer
func (s *TrainerService) GetTrainerAssignments(trainerID uint, status string) ([]TrainerAssignmentResponse, error) {
	var assignments []models.TrainerAssignment
	query := s.db.Where("trainer_id = ?", trainerID).
		Preload("Member").
		Preload("Plan").
		Preload("UserPlan").
		Preload("AssignedByUser")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("assigned_at DESC").Find(&assignments).Error; err != nil {
		return nil, err
	}

	var responses []TrainerAssignmentResponse
	for _, assignment := range assignments {
		responses = append(responses, *s.buildTrainerAssignmentResponse(&assignment))
	}

	return responses, nil
}

// GetAllTrainerAssignments gets all trainer assignments (admin only)
func (s *TrainerService) GetAllTrainerAssignments(status string) ([]TrainerAssignmentResponse, error) {
	var assignments []models.TrainerAssignment
	query := s.db.Preload("Trainer").
		Preload("Member").
		Preload("Plan").
		Preload("UserPlan").
		Preload("AssignedByUser")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("assigned_at DESC").Find(&assignments).Error; err != nil {
		return nil, err
	}

	var responses []TrainerAssignmentResponse
	for _, assignment := range assignments {
		responses = append(responses, *s.buildTrainerAssignmentResponse(&assignment))
	}

	return responses, nil
}

// GetMemberAssignments gets all trainer assignments for a specific member
func (s *TrainerService) GetMemberAssignments(memberID uint, status string) ([]TrainerAssignmentResponse, error) {
	var assignments []models.TrainerAssignment
	query := s.db.Where("member_id = ?", memberID).
		Preload("Trainer").
		Preload("Plan").
		Preload("UserPlan").
		Preload("AssignedByUser")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("assigned_at DESC").Find(&assignments).Error; err != nil {
		return nil, err
	}

	var responses []TrainerAssignmentResponse
	for _, assignment := range assignments {
		responses = append(responses, *s.buildTrainerAssignmentResponse(&assignment))
	}

	return responses, nil
}

// UpdateAssignmentStatus updates the status of a trainer assignment
func (s *TrainerService) UpdateAssignmentStatus(assignmentID uint, status string, trainerID uint) (*TrainerAssignmentResponse, error) {
	var assignment models.TrainerAssignment
	if err := s.db.First(&assignment, assignmentID).Error; err != nil {
		return nil, err
	}

	// Verify trainer owns this assignment
	if assignment.TrainerID != trainerID {
		return nil, gorm.ErrRecordNotFound
	}

	// Update status
	assignment.Status = status
	if err := s.db.Save(&assignment).Error; err != nil {
		return nil, err
	}

	// Also update the corresponding UserPlan status
	var userPlan models.UserPlan
	if err := s.db.First(&userPlan, assignment.UserPlanID).Error; err == nil {
		userPlan.Status = status
		s.db.Save(&userPlan)
	}

	// Reload with relationships
	if err := s.db.Preload("Trainer").
		Preload("Member").
		Preload("Plan").
		Preload("UserPlan").
		Preload("AssignedByUser").
		First(&assignment, assignment.ID).Error; err != nil {
		return nil, err
	}

	return s.buildTrainerAssignmentResponse(&assignment), nil
}

// GetAssignmentDetails gets detailed information about a specific assignment
func (s *TrainerService) GetAssignmentDetails(assignmentID uint, trainerID uint) (*TrainerAssignmentResponse, error) {
	var assignment models.TrainerAssignment
	if err := s.db.Where("id = ? AND trainer_id = ?", assignmentID, trainerID).
		Preload("Trainer").
		Preload("Member").
		Preload("Plan").
		Preload("UserPlan").
		Preload("AssignedByUser").
		First(&assignment).Error; err != nil {
		return nil, err
	}

	return s.buildTrainerAssignmentResponse(&assignment), nil
}

// Helper function to build response
func (s *TrainerService) buildTrainerAssignmentResponse(assignment *models.TrainerAssignment) *TrainerAssignmentResponse {
	response := &TrainerAssignmentResponse{
		ID:         assignment.ID,
		TrainerID:  assignment.TrainerID,
		MemberID:   assignment.MemberID,
		PlanID:     assignment.PlanID,
		UserPlanID: assignment.UserPlanID,
		Status:     assignment.Status,
		AssignedAt: assignment.AssignedAt.Format("2006-01-02 15:04:05"),
		AssignedBy: assignment.AssignedBy,
		Notes:      assignment.Notes,
	}

	// Add trainer info
	if assignment.Trainer.ID != 0 {
		response.Trainer = &UserInfo{
			ID:        assignment.Trainer.ID,
			FirstName: assignment.Trainer.FirstName,
			LastName:  assignment.Trainer.LastName,
			Email:     assignment.Trainer.Email,
			Role:      assignment.Trainer.Role,
		}
	}

	// Add member info
	if assignment.Member.ID != 0 {
		response.Member = &UserInfo{
			ID:        assignment.Member.ID,
			FirstName: assignment.Member.FirstName,
			LastName:  assignment.Member.LastName,
			Email:     assignment.Member.Email,
			Role:      assignment.Member.Role,
		}
	}

	// Add plan info
	if assignment.Plan.ID != 0 {
		response.Plan = &PlanResponse{
			ID:          assignment.Plan.ID,
			Name:        assignment.Plan.Name,
			Description: assignment.Plan.Description,
			GoalType:    assignment.Plan.GoalType,
			PlanType:    assignment.Plan.PlanType,
			Duration:    assignment.Plan.Duration,
			IsActive:    assignment.Plan.IsActive,
			CreatedAt:   assignment.Plan.CreatedAt,
			UpdatedAt:   assignment.Plan.UpdatedAt,
		}
	}

	// Add assigned by info
	if assignment.AssignedByUser.ID != 0 {
		response.AssignedByUser = &UserInfo{
			ID:        assignment.AssignedByUser.ID,
			FirstName: assignment.AssignedByUser.FirstName,
			LastName:  assignment.AssignedByUser.LastName,
			Email:     assignment.AssignedByUser.Email,
			Role:      assignment.AssignedByUser.Role,
		}
	}

	return response
}
