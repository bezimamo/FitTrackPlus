package plan

import (
	"net/http"
	"strconv"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// PlanHandler handles plan HTTP requests
type PlanHandler struct {
	planService        *PlanService
	planRequestService *PlanRequestService
	trainerService     *TrainerService
}

// NewPlanHandler creates a new plan handler
func NewPlanHandler(cfg *config.Config) *PlanHandler {
	return &PlanHandler{
		planService:        NewPlanService(cfg),
		planRequestService: NewPlanRequestService(cfg),
		trainerService:     NewTrainerService(cfg),
	}
}

// CreatePlan godoc
// @Summary Create a new plan template
// @Description Create a new fitness plan template (Admin/Trainer only)
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan body PlanRequest true "Plan details"
// @Success 201 {object} PlanResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin/Trainer only"
// @Router /plans [post]
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	// Check if user has permission (Admin or Trainer)
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found",
		})
		return
	}

	if userRole != "admin" && userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins and trainers can create plans",
		})
		return
	}

	// Get current user ID
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	// Bind request
	var req PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create plan
	plan, err := h.planService.CreatePlan(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// GetPlans godoc
// @Summary Get all plans
// @Description Get all plan templates with optional filtering
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goal_type query string false "Filter by goal type"
// @Param plan_type query string false "Filter by plan type"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {array} PlanResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /plans [get]
func (h *PlanHandler) GetPlans(c *gin.Context) {
	// Get query parameters
	goalType := c.Query("goal_type")
	planType := c.Query("plan_type")
	isActiveStr := c.Query("is_active")

	var isActive *bool
	if isActiveStr != "" {
		active, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &active
		}
	}

	// Get plans
	plans, err := h.planService.GetPlans(goalType, planType, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get plans",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetPlan godoc
// @Summary Get a specific plan
// @Description Get a specific plan template by ID
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} PlanResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Plan not found"
// @Router /plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	// Get plan ID from URL
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	// Get plan
	plan, err := h.planService.GetPlan(uint(planID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Plan not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// AssignPlan godoc
// @Summary Assign a plan to a user
// @Description Assign a plan template to a specific user (Trainer/Admin only)
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param assignment body map[string]interface{} true "Assignment details"
// @Success 201 {object} UserPlanResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer/Admin only"
// @Router /plans/assign [post]
func (h *PlanHandler) AssignPlan(c *gin.Context) {
	// Check if user has permission (Admin or Trainer)
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found",
		})
		return
	}

	if userRole != "admin" && userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins and trainers can assign plans",
		})
		return
	}

	// Get current user ID (who is assigning the plan)
	assignedBy, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	// Bind request
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
		PlanID uint `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Assign plan
	userPlan, err := h.planService.AssignPlan(req.UserID, req.PlanID, assignedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to assign plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, userPlan)
}

// GetUserPlans godoc
// @Summary Get user's assigned plans
// @Description Get all plans assigned to the current user
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} UserPlanResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /plans/my-plans [get]
func (h *PlanHandler) GetUserPlans(c *gin.Context) {
	// Get current user ID
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	// Get user's plans
	userPlans, err := h.planService.GetUserPlans(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user plans",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userPlans)
}

// GetAssignedPlans godoc
// @Summary Get plans assigned by trainer/admin
// @Description Get all plans assigned by the current trainer/admin (Trainer/Admin only)
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} UserPlanResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer/Admin only"
// @Router /plans/assigned [get]
func (h *PlanHandler) GetAssignedPlans(c *gin.Context) {
	// Check if user has permission (Admin or Trainer)
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found",
		})
		return
	}

	if userRole != "admin" && userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins and trainers can view assigned plans",
		})
		return
	}

	// For now, return all assigned plans
	// In a real app, you'd filter by the current trainer/admin
	c.JSON(http.StatusOK, gin.H{
		"message": "Assigned plans feature coming soon",
		"plans": []UserPlanResponse{},
	})
}

// GetAvailablePlans godoc
// @Summary Get available plans for members
// @Description Get all active plan templates that members can browse and select
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goal_type query string false "Filter by goal type"
// @Param plan_type query string false "Filter by plan type"
// @Success 200 {array} PlanResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /plans/available [get]
func (h *PlanHandler) GetAvailablePlans(c *gin.Context) {
	// Get query parameters
	goalType := c.Query("goal_type")
	planType := c.Query("plan_type")
	
	// Only show active plans to members
	isActive := true

	// Get available plans
	plans, err := h.planService.GetPlans(goalType, planType, &isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get available plans",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Available plans for selection",
		"plans": plans,
		"total": len(plans),
	})
}

// RequestPlanAssignment godoc
// @Summary Request plan assignment (Member only)
// @Description Request to be assigned a specific plan (Creates pending request for admin approval)
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePlanRequestRequest true "Plan request details"
// @Success 201 {object} PlanRequestResponse "Request created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Members only"
// @Router /plans/request [post]
func (h *PlanHandler) RequestPlanAssignment(c *gin.Context) {
	// Check if user is a member
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found",
		})
		return
	}

	if userRole != "member" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only members can request plan assignments",
		})
		return
	}

	// Get current user ID
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	// Bind request
	var req CreatePlanRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create the plan request
	planRequest, err := h.planRequestService.CreatePlanRequest(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create plan request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plan request submitted successfully! Waiting for admin approval.",
		"request": planRequest,
	})
}

// ===== ADMIN ENDPOINTS FOR PLAN REQUEST MANAGEMENT =====

// GetPendingRequests godoc
// @Summary Get all pending plan requests (Admin only)
// @Description Get all pending plan requests that need approval
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} PlanRequestResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/plan-requests/pending [get]
func (h *PlanHandler) GetPendingRequests(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view plan requests"})
		return
	}

	requests, err := h.planRequestService.GetPendingRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pending requests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pending plan requests",
		"requests": requests,
		"total": len(requests),
	})
}

// GetAllRequests godoc
// @Summary Get all plan requests (Admin only)
// @Description Get all plan requests with optional status filter
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (pending, approved, rejected)"
// @Success 200 {array} PlanRequestResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/plan-requests [get]
func (h *PlanHandler) GetAllRequests(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view plan requests"})
		return
	}

	status := c.Query("status")
	requests, err := h.planRequestService.GetAllRequests(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get requests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan requests",
		"requests": requests,
		"total": len(requests),
		"filter": status,
	})
}

// ApproveRequest godoc
// @Summary Approve a plan request and assign trainer (Admin only)
// @Description Approve a pending plan request and assign a trainer to manage it
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Request ID"
// @Param approval body ApprovalRequest true "Approval details with trainer assignment"
// @Success 200 {object} PlanRequestResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/plan-requests/{id}/approve [post]
func (h *PlanHandler) ApproveRequest(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can approve requests"})
		return
	}

	// Get admin user ID
	adminID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get request ID from URL
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Bind approval request
	var approval ApprovalRequest
	if err := c.ShouldBindJSON(&approval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid approval data",
			"details": err.Error(),
		})
		return
	}

	// Approve the request
	approvedRequest, err := h.planRequestService.ApproveRequest(uint(requestID), adminID, &approval)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to approve request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan request approved and trainer assigned successfully!",
		"request": approvedRequest,
	})
}

// RejectRequest godoc
// @Summary Reject a plan request (Admin only)
// @Description Reject a pending plan request
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Request ID"
// @Success 200 {object} PlanRequestResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Router /admin/plan-requests/{id}/reject [post]
func (h *PlanHandler) RejectRequest(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can reject requests"})
		return
	}

	// Get admin user ID
	adminID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get request ID from URL
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Reject the request
	rejectedRequest, err := h.planRequestService.RejectRequest(uint(requestID), adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to reject request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan request rejected",
		"request": rejectedRequest,
	})
}

// ===== MEMBER ENDPOINTS =====

// GetMyRequests godoc
// @Summary Get user's plan requests
// @Description Get all plan requests made by the current user
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} PlanRequestResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /plans/my-requests [get]
func (h *PlanHandler) GetMyRequests(c *gin.Context) {
	// Get current user ID
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	requests, err := h.planRequestService.GetUserRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get your requests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your plan requests",
		"requests": requests,
		"total": len(requests),
	})
}

// GetPlanRequestDetails godoc
// @Summary Get plan request details with trainer assignment (Admin only)
// @Description Get detailed information about a specific plan request including trainer assignment if approved
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan Request ID"
// @Success 200 {object} PlanRequestDetailResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Router /admin/plan-requests/{id} [get]
func (h *PlanHandler) GetPlanRequestDetails(c *gin.Context) {
	// Check if user is admin
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view request details"})
		return
	}

	// Get request ID from URL
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Get request details with trainer assignment
	details, err := h.planRequestService.GetRequestDetails(uint(requestID))
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get request details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request details retrieved successfully",
		"request": details,
	})
}

// ===== TRAINER ENDPOINTS =====

// GetMyAssignments godoc
// @Summary Get trainer's assigned members (Trainer only)
// @Description Get all members and plans assigned to the current trainer
// @Tags Trainer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (active, completed, paused)"
// @Success 200 {array} TrainerAssignmentResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer only"
// @Router /trainer/assignments [get]
func (h *PlanHandler) GetMyAssignments(c *gin.Context) {
	// Check if user is trainer
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can view assignments"})
		return
	}

	// Get trainer ID
	trainerID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	status := c.Query("status")
	assignments, err := h.trainerService.GetTrainerAssignments(trainerID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get assignments",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your assigned members and plans",
		"assignments": assignments,
		"total": len(assignments),
		"filter": status,
	})
}

// UpdateAssignmentStatus godoc
// @Summary Update assignment status (Trainer only)
// @Description Update the status of a trainer assignment
// @Tags Trainer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assignment ID"
// @Param update body map[string]string true "Status update"
// @Success 200 {object} TrainerAssignmentResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Trainer only"
// @Router /trainer/assignments/{id}/status [put]
func (h *PlanHandler) UpdateAssignmentStatus(c *gin.Context) {
	// Check if user is trainer
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if userRole != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can update assignments"})
		return
	}

	// Get trainer ID
	trainerID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get assignment ID from URL
	assignmentIDStr := c.Param("id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Bind status update
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status data",
			"details": err.Error(),
		})
		return
	}

	// Update assignment status
	assignment, err := h.trainerService.UpdateAssignmentStatus(uint(assignmentID), req.Status, trainerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update assignment status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment status updated successfully",
		"assignment": assignment,
	})
}
