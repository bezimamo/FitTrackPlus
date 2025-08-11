package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Test data structures
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	baseURL := "http://localhost:8080/api/v1"
	
	fmt.Println("🧪 Testing FitTrack+ Authentication System...")
	fmt.Println("=" * 50)

	// Test 1: Register a new user
	fmt.Println("\n1️⃣ Testing User Registration...")
	registerReq := RegisterRequest{
		Email:     "test@fittrackplus.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+1234567890",
		Role:      "member",
	}

	registerResp, err := makeRequest("POST", baseURL+"/auth/register", registerReq)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}

	var registerResult AuthResponse
	if err := json.Unmarshal(registerResp, &registerResult); err != nil {
		fmt.Printf("❌ Failed to parse registration response: %v\n", err)
		return
	}

	fmt.Printf("✅ Registration successful! User ID: %d\n", registerResult.User.ID)
	fmt.Printf("📧 Email: %s\n", registerResult.User.Email)
	fmt.Printf("🔑 Token: %s...\n", registerResult.Token[:20])

	// Test 2: Login with the registered user
	fmt.Println("\n2️⃣ Testing User Login...")
	loginReq := LoginRequest{
		Email:    "test@fittrackplus.com",
		Password: "password123",
	}

	loginResp, err := makeRequest("POST", baseURL+"/auth/login", loginReq)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return
	}

	var loginResult AuthResponse
	if err := json.Unmarshal(loginResp, &loginResult); err != nil {
		fmt.Printf("❌ Failed to parse login response: %v\n", err)
		return
	}

	fmt.Printf("✅ Login successful! User ID: %d\n", loginResult.User.ID)
	fmt.Printf("🔑 Token: %s...\n", loginResult.Token[:20])

	// Test 3: Get user profile (protected endpoint)
	fmt.Println("\n3️⃣ Testing Protected Profile Endpoint...")
	profileResp, err := makeAuthenticatedRequest("GET", baseURL+"/users/profile", nil, loginResult.Token)
	if err != nil {
		fmt.Printf("❌ Profile fetch failed: %v\n", err)
		return
	}

	var profileResult User
	if err := json.Unmarshal(profileResp, &profileResult); err != nil {
		fmt.Printf("❌ Failed to parse profile response: %v\n", err)
		return
	}

	fmt.Printf("✅ Profile fetch successful!\n")
	fmt.Printf("👤 Name: %s %s\n", profileResult.FirstName, profileResult.LastName)
	fmt.Printf("📧 Email: %s\n", profileResult.Email)
	fmt.Printf("📱 Phone: %s\n", profileResult.Phone)

	// Test 4: Test invalid login
	fmt.Println("\n4️⃣ Testing Invalid Login...")
	invalidLoginReq := LoginRequest{
		Email:    "test@fittrackplus.com",
		Password: "wrongpassword",
	}

	_, err = makeRequest("POST", baseURL+"/auth/login", invalidLoginReq)
	if err != nil {
		fmt.Printf("✅ Invalid login correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid login should have been rejected!")
	}

	fmt.Println("\n🎉 All authentication tests completed!")
}

func makeRequest(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func makeAuthenticatedRequest(method, url string, body interface{}, token string) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
} 