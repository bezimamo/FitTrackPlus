package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"fittrackplus/internal/auth"
	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/dashboard"
	"fittrackplus/internal/exercise"
	"fittrackplus/internal/plan"
	"fittrackplus/internal/profile"
	"fittrackplus/internal/workout"
	_ "fittrackplus/docs" // This is required for swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title FitTrack+ API
// @version 1.0
// @description A comprehensive fitness platform API for gym members and trainers
// @termsOfService http://swagger.io/terms/

// @contact.name FitTrack+ Team
// @contact.url http://www.fittrackplus.com/support
// @contact.email support@fittrackplus.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// main is the entry point of our application
// In Go, the main function must be in the main package
func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Connect to the database
	err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	// Gin is a popular HTTP web framework for Go
	router := gin.Default()

	// Add middleware for CORS (Cross-Origin Resource Sharing)
	// This allows our frontend to communicate with the backend
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Define our API routes
	// In Go, we use handlers (functions) to process HTTP requests
	setupRoutes(router, cfg)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("🚀 FitTrack+ Server starting on port %s\n", cfg.Port)
	fmt.Println("📊 API Documentation: http://localhost:" + cfg.Port + "/api/docs")
	fmt.Println("📚 Swagger UI: http://localhost:" + cfg.Port + "/swagger/index.html")
	fmt.Println("🔗 Health Check: http://localhost:" + cfg.Port + "/health")
	
	// ListenAndServe starts the HTTP server
	// If there's an error, log.Fatal will print it and exit
	log.Fatal(router.Run(":" + cfg.Port))
}

// setupRoutes defines all our API endpoints
// In Go, we group related functionality into functions
func setupRoutes(router *gin.Engine, cfg *config.Config) {
	// Create handlers
	authHandler := auth.NewAuthHandler(cfg)
	profileHandler := profile.NewProfileHandler(cfg)
	dashboardHandler := dashboard.NewDashboardHandler(cfg)
	planHandler := plan.NewPlanHandler(cfg)
	exerciseHandler := exercise.NewExerciseHandler(cfg)
	
	// Create workout service and handler
	workoutService := workout.NewWorkoutService(database.DB)
	workoutHandler := workout.NewWorkoutHandler(workoutService)
	
	// Create user workout service and handler
	userWorkoutService := workout.NewUserWorkoutService(database.DB)
	userWorkoutHandler := workout.NewUserWorkoutHandler(userWorkoutService)

	// Debug: Check if handlers are created successfully
	fmt.Println("🔧 Handlers initialized:")
	fmt.Println("   - AuthHandler:", authHandler != nil)
	fmt.Println("   - DashboardHandler:", dashboardHandler != nil)
	fmt.Println("   - PlanHandler:", planHandler != nil)
	fmt.Println("   - WorkoutHandler:", workoutHandler != nil)
	fmt.Println("   - ExerciseHandler:", exerciseHandler != nil)

	// API version 1 group
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", healthCheck)
		
		// Auth routes (public - no authentication required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}
		
		// User routes (protected - authentication required)
		users := api.Group("/users")
		users.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// Basic user profile (from auth)
			users.PUT("/profile", authHandler.UpdateProfile)
			
			// Enhanced profile management
			profileGroup := users.Group("/profile")
			{
				profileGroup.POST("/setup", profileHandler.SetupProfile)
				profileGroup.GET("", profileHandler.GetProfile)
				profileGroup.POST("/upload-image", profileHandler.UploadProfileImage)
				profileGroup.GET("/completion", profileHandler.CheckProfileCompletion)
				
				// Role-based profile management
				profileGroup.POST("/setup-role", profileHandler.SetupRoleProfile)
				profileGroup.GET("/role", profileHandler.GetRoleProfile)
				profileGroup.GET("/role/completion", profileHandler.CheckRoleProfileCompletion)
			}
		}

		// Dashboard routes (protected - authentication required)
		dashboardGroup := api.Group("/dashboard")
		dashboardGroup.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// General dashboard (automatically role-based)
			dashboardGroup.GET("", dashboardHandler.GetDashboard)
			
			// Dashboard components
			dashboardGroup.GET("/stats", dashboardHandler.GetDashboardStats)
			dashboardGroup.GET("/quick-actions", dashboardHandler.GetQuickActions)
			dashboardGroup.GET("/notifications", dashboardHandler.GetNotifications)
		}

		// Plan routes (protected - authentication required)
		planGroup := api.Group("/plans")
		planGroup.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// Plan management (Admin/Trainer only)
			planGroup.POST("", planHandler.CreatePlan)
			planGroup.GET("", planHandler.GetPlans)
			planGroup.GET("/:id", planHandler.GetPlan)
			
			// Plan assignment (Admin/Trainer only)
			planGroup.POST("/assign", planHandler.AssignPlan)
			
			// User plan access
			planGroup.GET("/my-plans", planHandler.GetUserPlans)
			planGroup.GET("/assigned", planHandler.GetAssignedPlans)
			
			// Member plan selection and requests
			planGroup.GET("/available", planHandler.GetAvailablePlans)
			planGroup.POST("/request", planHandler.RequestPlanAssignment)
			planGroup.GET("/my-requests", planHandler.GetMyRequests)
		}

		// Exercise routes (public - no authentication required for viewing)
		exerciseGroup := api.Group("/exercises")
		{
			// Public exercise access
			exerciseGroup.GET("", exerciseHandler.GetExercises)
			exerciseGroup.GET("/:id", exerciseHandler.GetExercise)
			exerciseGroup.GET("/categories", exerciseHandler.GetExerciseCategories)
			exerciseGroup.GET("/muscle-groups", exerciseHandler.GetMuscleGroups)
			exerciseGroup.GET("/difficulties", exerciseHandler.GetDifficulties)
			exerciseGroup.GET("/equipment", exerciseHandler.GetEquipment)
		}

		// Workout routes (public for viewing, protected for management)
		workoutGroup := api.Group("/workouts")
		{
			// Public workout access
			workoutGroup.GET("", workoutHandler.GetWorkouts)
			workoutGroup.GET("/:id", workoutHandler.GetWorkoutByID)
			workoutGroup.GET("/categories", workoutHandler.GetWorkoutCategories)
			workoutGroup.GET("/difficulties", workoutHandler.GetWorkoutDifficulties)
		}

		// User workout session routes (protected - authentication required)
		userWorkoutGroup := api.Group("/workouts")
		userWorkoutGroup.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// Workout session management
			userWorkoutGroup.POST("/start", userWorkoutHandler.StartWorkout)
			userWorkoutGroup.POST("/:id/complete-set", userWorkoutHandler.CompleteExerciseSet)
			userWorkoutGroup.POST("/:id/complete", userWorkoutHandler.CompleteWorkout)
			userWorkoutGroup.POST("/:id/pause", userWorkoutHandler.PauseWorkout)
			userWorkoutGroup.POST("/:id/resume", userWorkoutHandler.ResumeWorkout)
			userWorkoutGroup.GET("/session/:id", userWorkoutHandler.GetUserWorkout)
			userWorkoutGroup.GET("/my-sessions", userWorkoutHandler.GetUserWorkouts)
			userWorkoutGroup.GET("/:id/progress", userWorkoutHandler.GetWorkoutProgress)
		}

		// Admin routes for plan request management
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// Plan request management
			adminGroup.GET("/plan-requests", planHandler.GetAllRequests)
			adminGroup.GET("/plan-requests/pending", planHandler.GetPendingRequests)
			adminGroup.GET("/plan-requests/:id", planHandler.GetPlanRequestDetails)
			adminGroup.POST("/plan-requests/:id/approve", planHandler.ApproveRequest)
			adminGroup.POST("/plan-requests/:id/reject", planHandler.RejectRequest)
			
			// Exercise management (Admin only)
			adminGroup.POST("/exercises", exerciseHandler.CreateExercise)
			adminGroup.PUT("/exercises/:id", exerciseHandler.UpdateExercise)
			adminGroup.DELETE("/exercises/:id", exerciseHandler.DeleteExercise)
			adminGroup.GET("/exercises/stats", exerciseHandler.GetExerciseStats)
			// Media upload is now handled directly in CreateExercise endpoint
			
			// Workout management (Admin/Trainer only)
			adminGroup.POST("/workouts", workoutHandler.CreateWorkout)
			adminGroup.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
			adminGroup.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)
			adminGroup.GET("/workouts/stats", workoutHandler.GetWorkoutStats)
			
			                // User management (Admin only)
                adminGroup.GET("/users/trainers", authHandler.GetTrainers)
                adminGroup.GET("/users", authHandler.GetAllUsers)
                adminGroup.GET("/users/:id", authHandler.GetUserByID)
                adminGroup.POST("/users", authHandler.CreateUser)
                adminGroup.PUT("/users/:id", authHandler.UpdateUser)
                adminGroup.DELETE("/users/:id", authHandler.DeactivateUser)
                adminGroup.GET("/users/roles", authHandler.GetAvailableRoles)
                adminGroup.POST("/users/:id/change-role", authHandler.ChangeUserRole)
		}

		// Trainer routes for assignment management
		trainerGroup := api.Group("/trainer")
		trainerGroup.Use(auth.AuthMiddleware(cfg)) // Apply authentication middleware
		{
			// Assignment management
			trainerGroup.GET("/assignments", planHandler.GetMyAssignments)
			trainerGroup.PUT("/assignments/:id/status", planHandler.UpdateAssignmentStatus)
			
			// Workout management
			trainerGroup.GET("/workouts", workoutHandler.GetTrainerWorkouts)
		}
	}

	fmt.Println("✅ Routes configured successfully")
	fmt.Println("✅ - Auth routes: /api/v1/auth/*")
	fmt.Println("   - User routes: /api/v1/users/*")
	fmt.Println("   - Dashboard routes: /api/v1/dashboard/*")
	fmt.Println("   - Plan routes: /api/v1/plans/*")

	// Serve Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Serve API documentation
	router.GET("/api/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "FitTrack+ API Documentation",
			"version": "1.0.0",
			"swagger_ui": "http://localhost:" + cfg.Port + "/swagger/index.html",
			"endpoints": gin.H{
				"health": "/api/v1/health",
				"auth": gin.H{
					"register": "POST /api/v1/auth/register",
					"login": "POST /api/v1/auth/login",
				},
				"users": gin.H{
					"profile": "GET /api/v1/users/profile",
					"update": "PUT /api/v1/users/profile",
					"profile_setup": "POST /api/v1/users/profile/setup",
					"profile_image": "POST /api/v1/users/profile/upload-image",
					"profile_completion": "GET /api/v1/users/profile/completion",
					"role_profile_setup": "POST /api/v1/users/profile/setup-role",
					"role_profile": "GET /api/v1/users/profile/role",
					"role_profile_completion": "GET /api/v1/users/profile/role/completion",
				},
				"dashboard": gin.H{
					"general": "GET /api/v1/dashboard",
					"member": "GET /api/v1/dashboard/member",
					"trainer": "GET /api/v1/dashboard/trainer",
					"admin": "GET /api/v1/dashboard/admin",
					"stats": "GET /api/v1/dashboard/stats",
					"quick_actions": "GET /api/v1/dashboard/quick-actions",
					"notifications": "GET /api/v1/dashboard/notifications",
				},
				"plans": gin.H{
					"create": "POST /api/v1/plans",
					"list": "GET /api/v1/plans",
					"get": "GET /api/v1/plans/{id}",
					"assign": "POST /api/v1/plans/assign",
					"my_plans": "GET /api/v1/plans/my-plans",
					"assigned": "GET /api/v1/plans/assigned",
					"available": "GET /api/v1/plans/available",
					"request": "POST /api/v1/plans/request",
				},
			},
		})
	})
}

// Handler functions - these process HTTP requests
// Each handler receives a gin.Context which contains request/response data

// healthCheck is a simple endpoint to verify the server is running
// @Summary Health Check
// @Description Check if the API server is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "FitTrack+ API is running!",
		"timestamp": "2024-01-01T00:00:00Z", // We'll make this dynamic later
	})
} 