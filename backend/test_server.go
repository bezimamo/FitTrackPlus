package main

import (
	"fmt"
	"log"
	"net/http"

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

// This is a simple test to verify our basic setup works
// Run this with: go run test_server.go
func main() {
	// Create a simple router
	router := gin.Default()

	// Add a simple test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "FitTrack+ Test Server is running!",
			"status":  "success",
		})
	})

	// Add API routes for testing
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		// @Summary Health Check
		// @Description Check if the API server is running
		// @Tags Health
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Router /health [get]
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"message":   "FitTrack+ API is running!",
				"timestamp": "2024-01-01T00:00:00Z",
			})
		})

		// Auth routes (placeholders)
		auth := api.Group("/auth")
		{
			// @Summary Register User
			// @Description Register a new user account
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param user body map[string]interface{} true "User registration data"
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Router /auth/register [post]
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "User registration endpoint - coming soon!",
					"status":  "not_implemented",
				})
			})

			// @Summary Login User
			// @Description Authenticate user and return JWT token
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param credentials body map[string]interface{} true "User login credentials"
			// @Success 200 {object} map[string]interface{}
			// @Failure 401 {object} map[string]interface{}
			// @Router /auth/login [post]
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "User login endpoint - coming soon!",
					"status":  "not_implemented",
				})
			})
		}

		// User routes (placeholders)
		users := api.Group("/users")
		{
			// @Summary Get User Profile
			// @Description Get current user's profile information
			// @Tags Users
			// @Accept json
			// @Produce json
			// @Security BearerAuth
			// @Success 200 {object} map[string]interface{}
			// @Failure 401 {object} map[string]interface{}
			// @Router /users/profile [get]
			users.GET("/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Get user profile endpoint - coming soon!",
					"status":  "not_implemented",
				})
			})

			// @Summary Update User Profile
			// @Description Update current user's profile information
			// @Tags Users
			// @Accept json
			// @Produce json
			// @Security BearerAuth
			// @Param profile body map[string]interface{} true "Profile update data"
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Failure 401 {object} map[string]interface{}
			// @Router /users/profile [put]
			users.PUT("/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Update user profile endpoint - coming soon!",
					"status":  "not_implemented",
				})
			})
		}
	}

	// Serve Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	fmt.Println("🧪 Test server starting on port 8080")
	fmt.Println("🔗 Test endpoint: http://localhost:8080/test")
	fmt.Println("📚 Swagger UI: http://localhost:8080/swagger/index.html")
	fmt.Println("🔗 Health Check: http://localhost:8080/api/v1/health")
	
	log.Fatal(router.Run(":8080"))
} 