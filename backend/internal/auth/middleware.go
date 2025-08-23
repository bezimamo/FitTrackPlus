package auth

import (
	"fmt"
	"net/http"
	"strings"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/models"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates middleware for JWT authentication
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	authService := NewAuthService(cfg)

	return func(c *gin.Context) {
		fmt.Printf("🔐 Auth middleware called for: %s %s\n", c.Request.Method, c.Request.URL.Path)
		
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("📋 Authorization header: '%s'\n", authHeader)
		
		if authHeader == "" {
			fmt.Println("❌ No Authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Printf("❌ Authorization header doesn't start with 'Bearer ': '%s'\n", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("🔑 Token extracted: %s...\n", tokenString[:min(len(tokenString), 20)])

		// Validate the token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			fmt.Printf("❌ Token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		fmt.Printf("✅ Token validated for user ID: %d, role: %s\n", claims.UserID, claims.Role)

		// Get the user from database
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			fmt.Printf("❌ User not found: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			fmt.Printf("❌ User account deactivated: %d\n", claims.UserID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Account is deactivated",
			})
			c.Abort()
			return
		}

		// Store user information in context for later use
		c.Set("user", user)
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		fmt.Printf("✅ Authentication successful for user: %s (%s)\n", user.Email, user.Role)
		c.Next()
	}
}

// RoleMiddleware creates middleware to check user roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		allowed := false

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser extracts the current user from the context
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetCurrentUserID extracts the current user ID from the context
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserRole extracts the current user role from the context
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return userRole.(string), true
} 

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 