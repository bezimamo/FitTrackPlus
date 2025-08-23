package dashboard

import (
	"fmt"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// DashboardService handles dashboard logic
type DashboardService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(cfg *config.Config) *DashboardService {
	return &DashboardService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// DashboardResponse represents the dashboard data
type DashboardResponse struct {
	UserRole     string                 `json:"user_role"`
	UserInfo     UserInfo               `json:"user_info"`
	Stats        DashboardStats         `json:"stats"`
	RecentActivity []RecentActivity     `json:"recent_activity"`
	QuickActions []QuickAction          `json:"quick_actions"`
	Notifications []Notification        `json:"notifications"`
	// Role-specific data
	MemberData   *MemberDashboardData   `json:"member_data,omitempty"`
	TrainerData  *TrainerDashboardData  `json:"trainer_data,omitempty"`
	AdminData    *AdminDashboardData    `json:"admin_data,omitempty"`
}

// UserInfo contains basic user information
type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}

// DashboardStats contains general statistics
type DashboardStats struct {
	TotalUsers     int64 `json:"total_users"`
	ActivePlans    int64 `json:"active_plans"`
	TotalSessions  int64 `json:"total_sessions"`
	CompletionRate float64 `json:"completion_rate"`
}

// RecentActivity represents recent user activities
type RecentActivity struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"` // "login", "profile_update", "plan_assigned", etc.
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// QuickAction represents available quick actions
type QuickAction struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	Color       string `json:"color"`
}

// Notification represents user notifications
type Notification struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"` // "info", "warning", "success", "error"
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

// MemberDashboardData contains member-specific data
type MemberDashboardData struct {
	CurrentPlan     *PlanSummary     `json:"current_plan"`
	ProgressSummary ProgressSummary  `json:"progress_summary"`
	UpcomingSessions []SessionInfo   `json:"upcoming_sessions"`
	Goals           []GoalInfo       `json:"goals"`
}

// TrainerDashboardData contains trainer-specific data
type TrainerDashboardData struct {
	TotalClients    int64         `json:"total_clients"`
	ActiveClients   int64         `json:"active_clients"`
	TodaySessions   []SessionInfo `json:"today_sessions"`
	ClientProgress  []ClientProgress `json:"client_progress"`
}

// AdminDashboardData contains admin-specific data
type AdminDashboardData struct {
	UserStats       UserStats      `json:"user_stats"`
	RevenueStats    RevenueStats   `json:"revenue_stats"`
	SystemHealth    SystemHealth   `json:"system_health"`
	RecentSignups   []UserInfo     `json:"recent_signups"`
}

// PlanSummary represents a plan overview
type PlanSummary struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "workout", "diet", "mixed"
	Status      string    `json:"status"` // "active", "completed", "paused"
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Progress    float64   `json:"progress"` // 0-100
}

// ProgressSummary represents user progress
type ProgressSummary struct {
	CurrentWeight    float64 `json:"current_weight"`
	TargetWeight     float64 `json:"target_weight"`
	WeightLost       float64 `json:"weight_lost"`
	BodyFatReduction float64 `json:"body_fat_reduction"`
	MuscleGain       float64 `json:"muscle_gain"`
	OverallProgress  float64 `json:"overall_progress"`
}

// SessionInfo represents session information
type SessionInfo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"` // "workout", "consultation", "assessment"
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"` // minutes
	Status      string    `json:"status"` // "scheduled", "completed", "cancelled"
	TrainerName string    `json:"trainer_name,omitempty"`
	ClientName  string    `json:"client_name,omitempty"`
}

// GoalInfo represents user goals
type GoalInfo struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Target      float64 `json:"target"`
	Current     float64 `json:"current"`
	Progress    float64 `json:"progress"`
	Deadline    time.Time `json:"deadline"`
}

// ClientProgress represents client progress for trainers
type ClientProgress struct {
	ClientID    uint    `json:"client_id"`
	ClientName  string  `json:"client_name"`
	PlanName    string  `json:"plan_name"`
	Progress    float64 `json:"progress"`
	LastSession time.Time `json:"last_session"`
}

// UserStats represents user statistics for admin
type UserStats struct {
	TotalMembers  int64 `json:"total_members"`
	TotalTrainers int64 `json:"total_trainers"`
	ActiveUsers   int64 `json:"active_users"`
	NewThisMonth  int64 `json:"new_this_month"`
}

// RevenueStats represents revenue statistics for admin
type RevenueStats struct {
	MonthlyRevenue float64 `json:"monthly_revenue"`
	TotalRevenue   float64 `json:"total_revenue"`
	ActiveSubscriptions int `json:"active_subscriptions"`
	PendingPayments     int `json:"pending_payments"`
}

// SystemHealth represents system health for admin
type SystemHealth struct {
	DatabaseStatus string `json:"database_status"`
	APIServerStatus string `json:"api_server_status"`
	Uptime         string `json:"uptime"`
	LastBackup     time.Time `json:"last_backup"`
}

// GetDashboard retrieves dashboard data based on user role
func (s *DashboardService) GetDashboard(userID uint, userRole string) (*DashboardResponse, error) {
	// Get basic user info
	user, err := s.getUserInfo(userID)
	if err != nil {
		return nil, err
	}

	// Get general stats
	stats, err := s.getDashboardStats(userRole, userID)
	if err != nil {
		return nil, err
	}

	// Get recent activity
	recentActivity, err := s.getRecentActivity(userID)
	if err != nil {
		return nil, err
	}

	// Get quick actions based on role and profile completion
	quickActions := s.getQuickActions(userRole, userID)

	// Get notifications
	notifications, err := s.getNotifications(userID)
	if err != nil {
		return nil, err
	}

	response := &DashboardResponse{
		UserRole:        userRole,
		UserInfo:        *user,
		Stats:           *stats,
		RecentActivity:  recentActivity,
		QuickActions:    quickActions,
		Notifications:   notifications,
	}

	// Add role-specific data
	switch userRole {
	case "member":
		memberData, err := s.getMemberDashboardData(userID)
		if err != nil {
			return nil, err
		}
		response.MemberData = memberData
	case "trainer":
		trainerData, err := s.getTrainerDashboardData(userID)
		if err != nil {
			return nil, err
		}
		response.TrainerData = trainerData
	case "admin":
		adminData, err := s.getAdminDashboardData()
		if err != nil {
			return nil, err
		}
		response.AdminData = adminData
	}

	return response, nil
}

// calculateProfileCompletion calculates the completion percentage of a user profile
// This should match the calculation in the profile service
func (s *DashboardService) calculateProfileCompletion(profile *models.UserProfile) float64 {
	fields := 0
	total := 0

	// Basic info (required)
	if profile.Height > 0 {
		fields++
	}
	total++
	if profile.Weight > 0 {
		fields++
	}
	total++
	if profile.Age > 0 {
		fields++
	}
	total++
	if profile.Gender != "" {
		fields++
	}
	total++

	// Goals (required)
	if profile.Goals != "" {
		fields++
	}
	total++
	if profile.TargetWeight > 0 {
		fields++
	}
	total++
	if profile.Timeline > 0 {
		fields++
	}
	total++

	// Preferences (required)
	if profile.PreferredWorkoutTime != "" {
		fields++
	}
	total++
	if profile.WorkoutDays != "" {
		fields++
	}
	total++
	if profile.CommunicationPreference != "" {
		fields++
	}
	total++

	// Optional fields
	if profile.MedicalHistory != "" {
		fields++
	}
	total++
	if profile.ProfileImageURL != "" {
		fields++
	}
	total++

	return float64(fields) / float64(total) * 100
}

// getUserInfo retrieves basic user information
func (s *DashboardService) getUserInfo(userID uint) (*UserInfo, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}, nil
}

// getDashboardStats retrieves general dashboard statistics
func (s *DashboardService) getDashboardStats(userRole string, userID uint) (*DashboardStats, error) {
	var totalUsers, activePlans, totalSessions int64
	var completionRate float64

	// Get total users
	s.db.Model(&models.User{}).Count(&totalUsers)

	// Get active plans
	s.db.Model(&models.UserPlan{}).Where("status = ?", "active").Count(&activePlans)

	// Get total sessions
	s.db.Model(&models.Booking{}).Count(&totalSessions)

	// Calculate completion rate based on user role
	if userRole == "member" {
		// For members, show their personal profile completion
		var userProfile models.UserProfile
		err := s.db.Where("user_id = ?", userID).First(&userProfile).Error
		if err != nil {
			// Profile doesn't exist - return 0%
			completionRate = 0.0
			fmt.Printf("🔍 Member %d: No profile found, completion rate: 0%%\n", userID)
		} else {
			if userProfile.IsProfileComplete {
				completionRate = 100.0
				fmt.Printf("🔍 Member %d: Profile complete, completion rate: 100%%\n", userID)
			} else {
				// Calculate completion percentage based on filled fields
				completionRate = s.calculateProfileCompletion(&userProfile)
				fmt.Printf("🔍 Member %d: Profile incomplete, completion rate: %.1f%%\n", userID, completionRate)
			}
		}
	} else {
		// For admin/trainer, show system-wide completion rate
		if totalUsers > 0 {
			var completedProfiles int64
			s.db.Model(&models.UserProfile{}).Where("is_profile_complete = ?", true).Count(&completedProfiles)
			completionRate = float64(completedProfiles) / float64(totalUsers) * 100
		}
	}

	return &DashboardStats{
		TotalUsers:     totalUsers,
		ActivePlans:    activePlans,
		TotalSessions:  totalSessions,
		CompletionRate: completionRate,
	}, nil
}

// getRecentActivity retrieves recent user activities
func (s *DashboardService) getRecentActivity(userID uint) ([]RecentActivity, error) {
	// For now, return mock data
	// In a real app, you'd have an activities table
	activities := []RecentActivity{
		{
			ID:          1,
			Type:        "login",
			Description: "Successfully logged in",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          2,
			Type:        "profile_update",
			Description: "Updated profile information",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		},
	}

	return activities, nil
}

// getQuickActions returns quick actions based on user role and profile completion
func (s *DashboardService) getQuickActions(userRole string, userID uint) []QuickAction {
	switch userRole {
	case "member":
		// Check if user has completed their profile
		var userProfile models.UserProfile
		err := s.db.Where("user_id = ?", userID).First(&userProfile).Error
		
		if err != nil || !userProfile.IsProfileComplete {
			// Profile incomplete - show profile completion actions
			return []QuickAction{
				{
					ID:          "complete_profile",
					Title:       "Complete Profile",
					Description: "Complete your profile to get started",
					Icon:        "check",
					URL:         "/profile/setup",
					Color:       "orange",
				},
				{
					ID:          "update_profile",
					Title:       "Update Profile",
					Description: "Update your personal information",
					Icon:        "user",
					URL:         "/profile",
					Color:       "blue",
				},
			}
		}
		
		// Profile complete - show full actions
		return []QuickAction{
			{
				ID:          "update_profile",
				Title:       "Update Profile",
				Description: "Update your personal information",
				Icon:        "user",
				URL:         "/profile",
				Color:       "blue",
			},
			{
				ID:          "book_session",
				Title:       "Book Session",
				Description: "Schedule a training session",
				Icon:        "calendar",
				URL:         "/booking",
				Color:       "green",
			},
			{
				ID:          "log_progress",
				Title:       "Log Progress",
				Description: "Record your fitness progress",
				Icon:        "chart",
				URL:         "/progress",
				Color:       "purple",
			},
		}
	case "trainer":
		return []QuickAction{
			{
				ID:          "view_clients",
				Title:       "View Clients",
				Description: "Manage your client list",
				Icon:        "users",
				URL:         "/clients",
				Color:       "blue",
			},
			{
				ID:          "create_plan",
				Title:       "Create Plan",
				Description: "Create a new workout plan",
				Icon:        "plus",
				URL:         "/plans/create",
				Color:       "green",
			},
			{
				ID:          "schedule",
				Title:       "View Schedule",
				Description: "Check your session schedule",
				Icon:        "calendar",
				URL:         "/schedule",
				Color:       "orange",
			},
		}
	case "admin":
		return []QuickAction{
			{
				ID:          "manage_users",
				Title:       "Manage Users",
				Description: "View and manage all users",
				Icon:        "users",
				URL:         "/admin/users",
				Color:       "blue",
			},
			{
				ID:          "system_stats",
				Title:       "System Stats",
				Description: "View system analytics",
				Icon:        "chart",
				URL:         "/admin/stats",
				Color:       "green",
			},
			{
				ID:          "manage_plans",
				Title:       "Manage Plans",
				Description: "Create and manage plan templates",
				Icon:        "settings",
				URL:         "/admin/plans",
				Color:       "purple",
			},
		}
	default:
		return []QuickAction{}
	}
}

// getNotifications retrieves user notifications
func (s *DashboardService) getNotifications(userID uint) ([]Notification, error) {
	// Check if user has completed their profile
	var userProfile models.UserProfile
	err := s.db.Where("user_id = ?", userID).First(&userProfile).Error
	
	notifications := []Notification{}
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No profile found - show profile completion notification
			notifications = append(notifications, Notification{
				ID:        1,
				Type:      "warning",
				Title:     "Profile Incomplete",
				Message:   "Please complete your profile to access all features.",
				IsRead:    false,
				CreatedAt: time.Now().Add(-1 * time.Hour),
			})
		}
	} else if !userProfile.IsProfileComplete {
		// Profile exists but incomplete
		notifications = append(notifications, Notification{
			ID:        1,
			Type:      "warning",
			Title:     "Profile Incomplete",
			Message:   "Your profile is incomplete. Complete it to get personalized plans.",
			IsRead:    false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		})
	} else {
		// Profile complete - show welcome message
		notifications = append(notifications, Notification{
			ID:        1,
			Type:      "success",
			Title:     "Welcome!",
			Message:   "Your profile is complete. You can now access all features.",
			IsRead:    false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		})
	}

	return notifications, nil
}

// getMemberDashboardData retrieves member-specific dashboard data
func (s *DashboardService) getMemberDashboardData(userID uint) (*MemberDashboardData, error) {
	// First, check if user has completed their profile
	var userProfile models.UserProfile
	err := s.db.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No profile found - return empty data
			return &MemberDashboardData{
				CurrentPlan:      nil,
				ProgressSummary:  ProgressSummary{},
				UpcomingSessions: []SessionInfo{},
				Goals:           []GoalInfo{},
			}, nil
		}
		return nil, err
	}

	// Check if profile is complete (you can adjust this logic)
	if !userProfile.IsProfileComplete {
		// Profile incomplete - return empty data
		return &MemberDashboardData{
			CurrentPlan:      nil,
			ProgressSummary:  ProgressSummary{},
			UpcomingSessions: []SessionInfo{},
			Goals:           []GoalInfo{},
		}, nil
	}

	// Profile is complete - get real data
	// For now, still return mock data since we don't have real plans/sessions yet
	// In a real app, you'd query the actual database tables
	
	// Get current plan (mock for now)
	currentPlan, err := s.getCurrentPlan(userID)
	if err != nil {
		return nil, err
	}

	// Get progress summary (mock for now)
	progressSummary, err := s.getProgressSummary(userID)
	if err != nil {
		return nil, err
	}

	// Get upcoming sessions (mock for now)
	upcomingSessions, err := s.getUpcomingSessions(userID)
	if err != nil {
		return nil, err
	}

	// Get goals (mock for now)
	goals, err := s.getGoals(userID)
	if err != nil {
		return nil, err
	}

	return &MemberDashboardData{
		CurrentPlan:      currentPlan,
		ProgressSummary:  *progressSummary,
		UpcomingSessions: upcomingSessions,
		Goals:           goals,
	}, nil
}

// getTrainerDashboardData retrieves trainer-specific dashboard data
func (s *DashboardService) getTrainerDashboardData(userID uint) (*TrainerDashboardData, error) {
	// Get client statistics
	var totalClients, activeClients int64
	s.db.Model(&models.User{}).Where("role = ?", "member").Count(&totalClients)
	s.db.Model(&models.User{}).Where("role = ? AND is_active = ?", "member", true).Count(&activeClients)

	// Get today's sessions
	todaySessions, err := s.getTodaySessions(userID)
	if err != nil {
		return nil, err
	}

	// Get client progress
	clientProgress, err := s.getClientProgress(userID)
	if err != nil {
		return nil, err
	}

	return &TrainerDashboardData{
		TotalClients:   totalClients,
		ActiveClients:  activeClients,
		TodaySessions:  todaySessions,
		ClientProgress: clientProgress,
	}, nil
}

// getAdminDashboardData retrieves admin-specific dashboard data
func (s *DashboardService) getAdminDashboardData() (*AdminDashboardData, error) {
	// Get user statistics
	userStats, err := s.getUserStats()
	if err != nil {
		return nil, err
	}

	// Get revenue statistics
	revenueStats, err := s.getRevenueStats()
	if err != nil {
		return nil, err
	}

	// Get system health
	systemHealth, err := s.getSystemHealth()
	if err != nil {
		return nil, err
	}

	// Get recent signups
	recentSignups, err := s.getRecentSignups()
	if err != nil {
		return nil, err
	}

	return &AdminDashboardData{
		UserStats:     *userStats,
		RevenueStats:  *revenueStats,
		SystemHealth:  *systemHealth,
		RecentSignups: recentSignups,
	}, nil
}

// Helper methods for member dashboard
func (s *DashboardService) getCurrentPlan(userID uint) (*PlanSummary, error) {
	// For now, return mock data
	// In a real app, you'd query the user_plans table
	return &PlanSummary{
		ID:        1,
		Name:      "Beginner Fitness Plan",
		Type:      "workout",
		Status:    "active",
		StartDate: time.Now().AddDate(0, 0, -7),
		EndDate:   time.Now().AddDate(0, 1, 0),
		Progress:  25.0,
	}, nil
}

func (s *DashboardService) getProgressSummary(userID uint) (*ProgressSummary, error) {
	// For now, return mock data
	// In a real app, you'd calculate from progress_logs
	return &ProgressSummary{
		CurrentWeight:    75.0,
		TargetWeight:     70.0,
		WeightLost:       2.5,
		BodyFatReduction: 1.2,
		MuscleGain:       1.0,
		OverallProgress:  25.0,
	}, nil
}

func (s *DashboardService) getUpcomingSessions(userID uint) ([]SessionInfo, error) {
	// For now, return mock data
	// In a real app, you'd query the bookings table
	return []SessionInfo{
		{
			ID:          1,
			Title:       "Personal Training Session",
			Type:        "workout",
			Date:        time.Now().AddDate(0, 0, 1),
			Duration:    60,
			Status:      "scheduled",
			TrainerName: "John Trainer",
		},
	}, nil
}

func (s *DashboardService) getGoals(userID uint) ([]GoalInfo, error) {
	// For now, return mock data
	// In a real app, you'd query the goals table
	return []GoalInfo{
		{
			ID:          1,
			Title:       "Weight Loss",
			Description: "Lose 5kg in 3 months",
			Target:      70.0,
			Current:     75.0,
			Progress:    50.0,
			Deadline:    time.Now().AddDate(0, 2, 0),
		},
	}, nil
}

// Helper methods for trainer dashboard
func (s *DashboardService) getTodaySessions(userID uint) ([]SessionInfo, error) {
	// For now, return mock data
	return []SessionInfo{
		{
			ID:         1,
			Title:      "Morning Session with Alice",
			Type:       "workout",
			Date:       time.Now(),
			Duration:   60,
			Status:     "scheduled",
			ClientName: "Alice Johnson",
		},
	}, nil
}

func (s *DashboardService) getClientProgress(userID uint) ([]ClientProgress, error) {
	// For now, return mock data
	return []ClientProgress{
		{
			ClientID:    1,
			ClientName:  "Alice Johnson",
			PlanName:    "Weight Loss Plan",
			Progress:    75.0,
			LastSession: time.Now().AddDate(0, 0, -2),
		},
	}, nil
}

// Helper methods for admin dashboard
func (s *DashboardService) getUserStats() (*UserStats, error) {
	var totalMembers, totalTrainers, activeUsers, newThisMonth int64

	s.db.Model(&models.User{}).Where("role = ?", "member").Count(&totalMembers)
	s.db.Model(&models.User{}).Where("role = ?", "trainer").Count(&totalTrainers)
	s.db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
	s.db.Model(&models.User{}).Where("created_at >= ?", time.Now().AddDate(0, -1, 0)).Count(&newThisMonth)

	return &UserStats{
		TotalMembers:  totalMembers,
		TotalTrainers: totalTrainers,
		ActiveUsers:   activeUsers,
		NewThisMonth:  newThisMonth,
	}, nil
}

func (s *DashboardService) getRevenueStats() (*RevenueStats, error) {
	// For now, return mock data
	// In a real app, you'd calculate from payments table
	return &RevenueStats{
		MonthlyRevenue:      5000.0,
		TotalRevenue:        25000.0,
		ActiveSubscriptions: 50,
		PendingPayments:     5,
	}, nil
}

func (s *DashboardService) getSystemHealth() (*SystemHealth, error) {
	return &SystemHealth{
		DatabaseStatus: "healthy",
		APIServerStatus: "running",
		Uptime:         "7 days",
		LastBackup:     time.Now().AddDate(0, 0, -1),
	}, nil
}

func (s *DashboardService) getRecentSignups() ([]UserInfo, error) {
	var users []models.User
	s.db.Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).Limit(5).Find(&users)

	var userInfos []UserInfo
	for _, user := range users {
		userInfos = append(userInfos, UserInfo{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
		})
	}

	return userInfos, nil
}
