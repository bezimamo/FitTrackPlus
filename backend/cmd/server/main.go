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
	"fittrackplus/internal/plan"
	"fittrackplus/internal/profile"
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

	// Debug: Check if handlers are created successfully
	fmt.Println("🔧 Handlers initialized:")
	fmt.Println("   - AuthHandler:", authHandler != nil)
	fmt.Println("   - ProfileHandler:", profileHandler != nil)
	fmt.Println("   - DashboardHandler:", dashboardHandler != nil)
	fmt.Println("   - PlanHandler:", planHandler != nil)

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
			
			// Member plan selection
			planGroup.GET("/available", planHandler.GetAvailablePlans)
			planGroup.POST("/request", planHandler.RequestPlanAssignment)
		}
	}

	fmt.Println("✅ Routes configured successfully")
	fmt.Println("   - Auth routes: /api/v1/auth/*")
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