package dashboard

import (
	"fmt"
	"net/http"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard HTTP requests
type DashboardHandler struct {
	dashboardService *DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: NewDashboardService(cfg),
	}
}

// GetDashboard godoc
// @Summary Get user dashboard
// @Description Get personalized dashboard data based on user role
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DashboardResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	fmt.Println("🔍 Dashboard handler called!")
	
	// Get current user info from context
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		fmt.Println("❌ User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		fmt.Println("❌ User role not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	fmt.Printf("✅ User authenticated - ID: %d, Role: %s\n", userID, userRole)

	// Get dashboard data based on user role
	dashboardData, err := h.dashboardService.GetDashboard(userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get dashboard data",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dashboardData)
}



// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get general dashboard statistics
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DashboardStats
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get dashboard stats
	stats, err := h.dashboardService.getDashboardStats(userRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get dashboard statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetQuickActions godoc
// @Summary Get quick actions
// @Description Get available quick actions based on user role and profile completion
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} QuickAction
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /dashboard/quick-actions [get]
func (h *DashboardHandler) GetQuickActions(c *gin.Context) {
	userRole, exists := auth.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User role not found in context",
		})
		return
	}

	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get quick actions
	quickActions := h.dashboardService.getQuickActions(userRole, userID)

	c.JSON(http.StatusOK, quickActions)
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get user notifications
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Notification
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dashboard/notifications [get]
func (h *DashboardHandler) GetNotifications(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get notifications
	notifications, err := h.dashboardService.getNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get notifications",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
