package main

import (
	"fmt"
	"log"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	
	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	db := database.GetDB()
	
	// Test completion rate for different users
	userIDs := []uint{1, 2, 3, 4} // Test all users
	
	for _, userID := range userIDs {
		// Get user info
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			fmt.Printf("❌ User %d: Not found\n", userID)
			continue
		}
		
		fmt.Printf("\n👤 User %d: %s %s (%s)\n", userID, user.FirstName, user.LastName, user.Role)
		
		// Get profile
		var profile models.UserProfile
		if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			fmt.Printf("   📊 Profile: Not found (0%%)\n")
		} else {
			fmt.Printf("   📊 Profile: Found (ID: %d)\n", profile.ID)
			fmt.Printf("   ✅ Is Complete: %t\n", profile.IsProfileComplete)
			
			// Calculate completion manually
			completion := calculateCompletion(&profile)
			fmt.Printf("   📈 Completion: %.1f%%\n", completion)
		}
	}
	
	// System-wide stats
	var totalUsers, completedProfiles int64
	db.Model(&models.User{}).Count(&totalUsers)
	db.Model(&models.UserProfile{}).Where("is_profile_complete = ?", true).Count(&completedProfiles)
	
	fmt.Printf("\n📊 System-wide Stats:\n")
	fmt.Printf("   👥 Total Users: %d\n", totalUsers)
	fmt.Printf("   ✅ Completed Profiles: %d\n", completedProfiles)
	fmt.Printf("   📈 System Completion Rate: %.1f%%\n", float64(completedProfiles)/float64(totalUsers)*100)
}

func calculateCompletion(profile *models.UserProfile) float64 {
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

