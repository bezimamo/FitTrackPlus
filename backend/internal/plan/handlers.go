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
	planService *PlanService
}

// NewPlanHandler creates a new plan handler
func NewPlanHandler(cfg *config.Config) *PlanHandler {
	return &PlanHandler{
		planService: NewPlanService(cfg),
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
// @Description Request to be assigned a specific plan (Members can request, but admin/trainer must approve)
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Plan request details"
// @Success 201 {object} map[string]interface{} "Request submitted"
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
	var req struct {
		PlanID uint   `json:"plan_id" binding:"required"`
		Reason string `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if plan exists and is active
	_, err := h.planService.GetPlan(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Plan not found or not available",
		})
		return
	}

	// For now, directly assign the plan (in a real app, you'd create a request system)
	userPlan, err := h.planService.AssignPlan(userID, req.PlanID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to assign plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plan assigned successfully!",
		"plan": userPlan,
		"note": "In a real app, this would go through an approval process",
	})
}
