package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// ProfileService handles all profile-related operations
type ProfileService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewProfileService creates a new profile service
func NewProfileService(cfg *config.Config) *ProfileService {
	return &ProfileService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// ProfileSetupRequest represents the complete profile setup data
type ProfileSetupRequest struct {
	// Basic Information
	Height   float64 `json:"height" binding:"required,min=50,max=300"`
	Weight   float64 `json:"weight" binding:"required,min=20,max=500"`
	Age      int     `json:"age" binding:"required,min=13,max=120"`
	Gender   string  `json:"gender" binding:"required,oneof=male female other"`

	// Fitness Goals
	Goals        []string `json:"goals" binding:"required,min=1"`
	TargetWeight float64  `json:"target_weight" binding:"required,min=20,max=500"`
	Timeline     int      `json:"timeline" binding:"required,min=30,max=3650"`

	// Medical Information (optional)
	MedicalHistory string `json:"medical_history"`
	Allergies      string `json:"allergies"`
	Medications    string `json:"medications"`
	PhysioNeeds    string `json:"physio_needs"`

	// Physical Measurements (optional)
	BodyFatPercentage float64            `json:"body_fat_percentage"`
	MuscleMass        float64            `json:"muscle_mass"`
	BodyMeasurements  map[string]float64 `json:"body_measurements"`

	// Preferences
	PreferredWorkoutTime    string   `json:"preferred_workout_time" binding:"required"`
	WorkoutDays             []string `json:"workout_days" binding:"required,min=1"`
	CommunicationPreference string   `json:"communication_preference" binding:"required,oneof=email phone sms"`
}

// ProfileResponse represents the profile response
type ProfileResponse struct {
	UserProfile models.UserProfile `json:"profile"`
	IsComplete  bool               `json:"is_complete"`
	Completion  float64            `json:"completion_percentage"`
}

// SetupProfile creates or updates a complete user profile
func (s *ProfileService) SetupProfile(userID uint, req *ProfileSetupRequest) (*ProfileResponse, error) {
	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Convert arrays to JSON strings
	goalsJSON, err := json.Marshal(req.Goals)
	if err != nil {
		return nil, err
	}

	workoutDaysJSON, err := json.Marshal(req.WorkoutDays)
	if err != nil {
		return nil, err
	}

	bodyMeasurementsJSON, err := json.Marshal(req.BodyMeasurements)
	if err != nil {
		return nil, err
	}

	// Create or update profile
	var userProfile models.UserProfile
	result := s.db.Where("user_id = ?", userID).First(&userProfile)

	if result.Error != nil {
		// Profile doesn't exist, create new one
		userProfile = models.UserProfile{
			UserID:                  userID,
			Height:                  req.Height,
			Weight:                  req.Weight,
			Age:                     req.Age,
			Gender:                  req.Gender,
			Goals:                   string(goalsJSON),
			TargetWeight:            req.TargetWeight,
			Timeline:                req.Timeline,
			MedicalHistory:          req.MedicalHistory,
			Allergies:               req.Allergies,
			Medications:             req.Medications,
			PhysioNeeds:             req.PhysioNeeds,
			BodyFatPercentage:       req.BodyFatPercentage,
			MuscleMass:              req.MuscleMass,
			BodyMeasurements:        string(bodyMeasurementsJSON),
			PreferredWorkoutTime:    req.PreferredWorkoutTime,
			WorkoutDays:             string(workoutDaysJSON),
			CommunicationPreference: req.CommunicationPreference,
			IsProfileComplete:       true,
		}

		if err := s.db.Create(&userProfile).Error; err != nil {
			return nil, err
		}
	} else {
		// Profile exists, update it
		userProfile.Height = req.Height
		userProfile.Weight = req.Weight
		userProfile.Age = req.Age
		userProfile.Gender = req.Gender
		userProfile.Goals = string(goalsJSON)
		userProfile.TargetWeight = req.TargetWeight
		userProfile.Timeline = req.Timeline
		userProfile.MedicalHistory = req.MedicalHistory
		userProfile.Allergies = req.Allergies
		userProfile.Medications = req.Medications
		userProfile.PhysioNeeds = req.PhysioNeeds
		userProfile.BodyFatPercentage = req.BodyFatPercentage
		userProfile.MuscleMass = req.MuscleMass
		userProfile.BodyMeasurements = string(bodyMeasurementsJSON)
		userProfile.PreferredWorkoutTime = req.PreferredWorkoutTime
		userProfile.WorkoutDays = string(workoutDaysJSON)
		userProfile.CommunicationPreference = req.CommunicationPreference
		userProfile.IsProfileComplete = true

		if err := s.db.Save(&userProfile).Error; err != nil {
			return nil, err
		}
	}

	return s.buildProfileResponse(&userProfile), nil
}

// GetProfile retrieves a user's profile
func (s *ProfileService) GetProfile(userID uint) (*ProfileResponse, error) {
	var userProfile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Profile doesn't exist, return empty profile
			return &ProfileResponse{
				UserProfile: models.UserProfile{UserID: userID},
				IsComplete:  false,
				Completion:  0,
			}, nil
		}
		return nil, err
	}

	return s.buildProfileResponse(&userProfile), nil
}

// UploadProfileImage handles profile image upload
func (s *ProfileService) UploadProfileImage(userID uint, file *multipart.FileHeader) (*ProfileResponse, error) {
	// Validate file
	if err := s.validateImageFile(file); err != nil {
		return nil, err
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/profile_images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("profile_%d_%d%s", userID, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	if err := s.saveUploadedFile(file, filepath); err != nil {
		return nil, err
	}

	// Update profile with image URL
	var userProfile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
		// Create profile if it doesn't exist
		userProfile = models.UserProfile{
			UserID: userID,
		}
	}

	userProfile.ProfileImageURL = "/uploads/profile_images/" + filename
	completion := s.calculateProfileCompletion(&userProfile)
	userProfile.IsProfileComplete = completion >= 80.0 // Consider complete if 80% or more

	if err := s.db.Save(&userProfile).Error; err != nil {
		return nil, err
	}

	return s.buildProfileResponse(&userProfile), nil
}

// Helper methods
func (s *ProfileService) buildProfileResponse(profile *models.UserProfile) *ProfileResponse {
	completion := s.calculateProfileCompletion(profile)
	return &ProfileResponse{
		UserProfile: *profile,
		IsComplete:  profile.IsProfileComplete,
		Completion:  completion,
	}
}

func (s *ProfileService) calculateProfileCompletion(profile *models.UserProfile) float64 {
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

func (s *ProfileService) validateImageFile(file *multipart.FileHeader) error {
	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return errors.New("file size too large (max 5MB)")
	}

	// Check file type
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"}
	contentType := file.Header.Get("Content-Type")
	
	allowed := false
	for _, t := range allowedTypes {
		if contentType == t {
			allowed = true
			break
		}
	}

	if !allowed {
		return errors.New("invalid file type (only JPEG, PNG, GIF allowed)")
	}

	return nil
}

func (s *ProfileService) saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy file content
	_, err = out.ReadFrom(src)
	return err
} 