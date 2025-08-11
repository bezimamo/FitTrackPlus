package main

import (
	"fmt"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
)

func main() {
	fmt.Println("🔍 Testing FitTrack+ Database Connection...")
	
	// Load configuration
	cfg := config.LoadConfig()
	
	fmt.Printf("📊 Database Config:\n")
	fmt.Printf("   Host: %s\n", cfg.DBHost)
	fmt.Printf("   Port: %s\n", cfg.DBPort)
	fmt.Printf("   User: %s\n", cfg.DBUser)
	fmt.Printf("   Database: %s\n", cfg.DBName)
	fmt.Printf("   Password: %s\n", "***hidden***")
	
	// Test database connection
	fmt.Println("\n🔌 Attempting to connect to database...")
	err := database.Connect(cfg)
	if err != nil {
		fmt.Printf("❌ Database connection failed: %v\n", err)
		fmt.Println("\n💡 Troubleshooting tips:")
		fmt.Println("   1. Make sure PostgreSQL is running")
		fmt.Println("   2. Check if the database 'fittrackplus' exists")
		fmt.Println("   3. Verify the password in your .env file")
		fmt.Println("   4. Try: CREATE DATABASE fittrackplus;")
		return
	}
	
	fmt.Println("✅ Database connection successful!")
	fmt.Println("✅ Tables created successfully!")
	
	// Close the connection
	database.Close()
	fmt.Println("🔌 Database connection closed.")
} 