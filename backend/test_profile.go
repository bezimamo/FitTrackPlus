package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080/api/v1"

// Test data structures
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

type ProfileSetupRequest struct {
	Height                 float64            `json:"height"`
	Weight                 float64            `json:"weight"`
	Age                    int                `json:"age"`
	Gender                 string             `json:"gender"`
	Goals                  []string           `json:"goals"`
	TargetWeight           float64            `json:"target_weight"`
	Timeline               int                `json:"timeline"`
	MedicalHistory         string             `json:"medical_history"`
	Allergies              string             `json:"allergies"`
	Medications            string             `json:"medications"`
	PhysioNeeds            string             `json:"physio_needs"`
	BodyFatPercentage      float64            `json:"body_fat_percentage"`
	MuscleMass             float64            `json:"muscle_mass"`
	BodyMeasurements       map[string]float64 `json:"body_measurements"`
	PreferredWorkoutTime   string             `json:"preferred_workout_time"`
	WorkoutDays            []string           `json:"workout_days"`
	CommunicationPreference string            `json:"communication_preference"`
}

type ProfileResponse struct {
	Profile struct {
		ID                    uint    `json:"id"`
		UserID                uint    `json:"user_id"`
		Height                float64 `json:"height"`
		Weight                float64 `json:"weight"`
		Age                   int     `json:"age"`
		Gender                string  `json:"gender"`
		Goals                 string  `json:"goals"`
		TargetWeight          float64 `json:"target_weight"`
		Timeline              int     `json:"timeline"`
		MedicalHistory        string  `json:"medical_history"`
		Allergies             string  `json:"allergies"`
		Medications           string  `json:"medications"`
		PhysioNeeds           string  `json:"physio_needs"`
		BodyFatPercentage     float64 `json:"body_fat_percentage"`
		MuscleMass            float64 `json:"muscle_mass"`
		BodyMeasurements      string  `json:"body_measurements"`
		ProfileImageURL       string  `json:"profile_image_url"`
		PreferredWorkoutTime  string  `json:"preferred_workout_time"`
		WorkoutDays           string  `json:"workout_days"`
		CommunicationPreference string `json:"communication_preference"`
		IsProfileComplete     bool    `json:"is_profile_complete"`
	} `json:"profile"`
	IsComplete bool    `json:"is_complete"`
	Completion float64 `json:"completion_percentage"`
}

type CompletionResponse struct {
	IsComplete        bool    `json:"is_complete"`
	CompletionPercentage float64 `json:"completion_percentage"`
	ProfileExists     bool    `json:"profile_exists"`
}

func main() {
	fmt.Println("🧪 Testing FitTrack+ Profile Management System")
	fmt.Println("=============================================")

	// Test 1: Register a new user
	fmt.Println("\n1️⃣ Testing User Registration...")
	registerData := map[string]interface{}{
		"email":     "profiletest@example.com",
		"password":  "testpassword123",
		"first_name": "Profile",
		"last_name":  "Tester",
		"role":      "member",
	}

	registerResp, err := makeRequest("POST", "/auth/register", registerData, "")
	if err != nil {
		log.Printf("❌ Registration failed: %v", err)
		return
	}
	fmt.Printf("✅ Registration successful: %s\n", registerResp)

	// Test 2: Login to get token
	fmt.Println("\n2️⃣ Testing User Login...")
	loginData := LoginRequest{
		Email:    "profiletest@example.com",
		Password: "testpassword123",
	}

	loginResp, err := makeRequest("POST", "/auth/login", loginData, "")
	if err != nil {
		log.Printf("❌ Login failed: %v", err)
		return
	}

	var loginResult LoginResponse
	if err := json.Unmarshal([]byte(loginResp), &loginResult); err != nil {
		log.Printf("❌ Failed to parse login response: %v", err)
		return
	}

	token := loginResult.Token
	fmt.Printf("✅ Login successful, token: %s...\n", token[:20])

	// Test 3: Check initial profile completion
	fmt.Println("\n3️⃣ Testing Initial Profile Completion Check...")
	completionResp, err := makeAuthenticatedRequest("GET", "/users/profile/completion", nil, token)
	if err != nil {
		log.Printf("❌ Profile completion check failed: %v", err)
		return
	}

	var completionResult CompletionResponse
	if err := json.Unmarshal([]byte(completionResp), &completionResult); err != nil {
		log.Printf("❌ Failed to parse completion response: %v", err)
		return
	}

	fmt.Printf("✅ Profile completion: %.1f%% (Complete: %t, Exists: %t)\n", 
		completionResult.CompletionPercentage, 
		completionResult.IsComplete, 
		completionResult.ProfileExists)

	// Test 4: Setup complete profile
	fmt.Println("\n4️⃣ Testing Complete Profile Setup...")
	profileData := ProfileSetupRequest{
		Height:                 175.0,
		Weight:                 70.0,
		Age:                    25,
		Gender:                 "male",
		Goals:                  []string{"lose_weight", "gain_muscle", "improve_flexibility"},
		TargetWeight:           65.0,
		Timeline:               90,
		MedicalHistory:         "No major medical issues",
		Allergies:              "None",
		Medications:            "None",
		PhysioNeeds:            "Lower back exercises",
		BodyFatPercentage:      15.0,
		MuscleMass:             55.0,
		BodyMeasurements: map[string]float64{
			"chest":  95.0,
			"waist":  80.0,
			"arms":   30.0,
			"thighs": 55.0,
		},
		PreferredWorkoutTime:   "morning",
		WorkoutDays:            []string{"monday", "wednesday", "friday"},
		CommunicationPreference: "email",
	}

	profileResp, err := makeAuthenticatedRequest("POST", "/users/profile/setup", profileData, token)
	if err != nil {
		log.Printf("❌ Profile setup failed: %v", err)
		return
	}

	var profileResult ProfileResponse
	if err := json.Unmarshal([]byte(profileResp), &profileResult); err != nil {
		log.Printf("❌ Failed to parse profile response: %v", err)
		return
	}

	fmt.Printf("✅ Profile setup successful!\n")
	fmt.Printf("   - Profile ID: %d\n", profileResult.Profile.ID)
	fmt.Printf("   - Height: %.1f cm\n", profileResult.Profile.Height)
	fmt.Printf("   - Weight: %.1f kg\n", profileResult.Profile.Weight)
	fmt.Printf("   - Goals: %s\n", profileResult.Profile.Goals)
	fmt.Printf("   - Completion: %.1f%%\n", profileResult.Completion)
	fmt.Printf("   - Is Complete: %t\n", profileResult.IsComplete)

	// Test 5: Get profile
	fmt.Println("\n5️⃣ Testing Get Profile...")
	getProfileResp, err := makeAuthenticatedRequest("GET", "/users/profile", nil, token)
	if err != nil {
		log.Printf("❌ Get profile failed: %v", err)
		return
	}

	var getProfileResult ProfileResponse
	if err := json.Unmarshal([]byte(getProfileResp), &getProfileResult); err != nil {
		log.Printf("❌ Failed to parse get profile response: %v", err)
		return
	}

	fmt.Printf("✅ Get profile successful!\n")
	fmt.Printf("   - User ID: %d\n", getProfileResult.Profile.UserID)
	fmt.Printf("   - Age: %d\n", getProfileResult.Profile.Age)
	fmt.Printf("   - Gender: %s\n", getProfileResult.Profile.Gender)
	fmt.Printf("   - Target Weight: %.1f kg\n", getProfileResult.Profile.TargetWeight)
	fmt.Printf("   - Timeline: %d days\n", getProfileResult.Profile.Timeline)

	// Test 6: Check final profile completion
	fmt.Println("\n6️⃣ Testing Final Profile Completion Check...")
	finalCompletionResp, err := makeAuthenticatedRequest("GET", "/users/profile/completion", nil, token)
	if err != nil {
		log.Printf("❌ Final profile completion check failed: %v", err)
		return
	}

	var finalCompletionResult CompletionResponse
	if err := json.Unmarshal([]byte(finalCompletionResp), &finalCompletionResult); err != nil {
		log.Printf("❌ Failed to parse final completion response: %v", err)
		return
	}

	fmt.Printf("✅ Final profile completion: %.1f%% (Complete: %t, Exists: %t)\n", 
		finalCompletionResult.CompletionPercentage, 
		finalCompletionResult.IsComplete, 
		finalCompletionResult.ProfileExists)

	fmt.Println("\n🎉 All Profile Management Tests Completed Successfully!")
	fmt.Println("=====================================================")
}

// Helper functions
func makeRequest(method, endpoint string, data interface{}, token string) (string, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func makeAuthenticatedRequest(method, endpoint string, data interface{}, token string) (string, error) {
	return makeRequest(method, endpoint, data, token)
} 