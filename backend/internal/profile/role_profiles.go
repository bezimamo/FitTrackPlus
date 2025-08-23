package profile

import (
	"errors"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/database"
	"fittrackplus/internal/common/models"

	"gorm.io/gorm"
)

// RoleProfileService handles role-specific profile operations
type RoleProfileService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewRoleProfileService creates a new role profile service
func NewRoleProfileService(cfg *config.Config) *RoleProfileService {
	return &RoleProfileService{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// RoleProfileSetupRequest represents the profile setup request for different roles
type RoleProfileSetupRequest struct {
	// Common fields for all roles
	Phone string `json:"phone" binding:"required"`
	
	// Member-specific fields
	Height                float64 `json:"height,omitempty"`
	Weight                float64 `json:"weight,omitempty"`
	Age                   int     `json:"age,omitempty"`
	Gender                string  `json:"gender,omitempty"`
	Goals                 string  `json:"goals,omitempty"`
	TargetWeight          float64 `json:"target_weight,omitempty"`
	Timeline              int     `json:"timeline,omitempty"`
	MedicalHistory        string  `json:"medical_history,omitempty"`
	Allergies             string  `json:"allergies,omitempty"`
	Medications           string  `json:"medications,omitempty"`
	PhysioNeeds           string  `json:"physio_needs,omitempty"`
	BodyFatPercentage     float64 `json:"body_fat_percentage,omitempty"`
	MuscleMass            float64 `json:"muscle_mass,omitempty"`
	BodyMeasurements      string  `json:"body_measurements,omitempty"`
	PreferredWorkoutTime  string  `json:"preferred_workout_time,omitempty"`
	WorkoutDays           string  `json:"workout_days,omitempty"`
	CommunicationPreference string `json:"communication_preference,omitempty"`

	// Trainer-specific fields
	Certifications        string `json:"certifications,omitempty"`
	Experience            int    `json:"experience,omitempty"`
	Specializations       string `json:"specializations,omitempty"`
	Bio                   string `json:"bio,omitempty"`
	Philosophy            string `json:"philosophy,omitempty"`
	Availability          string `json:"availability,omitempty"`
	SessionRates          string `json:"session_rates,omitempty"`
	PackageRates          string `json:"package_rates,omitempty"`
	Languages             string `json:"languages,omitempty"`
	Education             string `json:"education,omitempty"`
	Awards                string `json:"awards,omitempty"`
	PortfolioImages       string `json:"portfolio_images,omitempty"`
	BeforeAfterPhotos     string `json:"before_after_photos,omitempty"`
	Testimonials          string `json:"testimonials,omitempty"`
	PreferredContactMethod string `json:"preferred_contact_method,omitempty"`
	ResponseTime          string `json:"response_time,omitempty"`

	// Admin-specific fields
	AdminRole             string `json:"admin_role,omitempty"`
	AccessLevel           string `json:"access_level,omitempty"`
	Department            string `json:"department,omitempty"`
	EmergencyContact      string `json:"emergency_contact,omitempty"`
	EmergencyPhone        string `json:"emergency_phone,omitempty"`
	OfficeLocation        string `json:"office_location,omitempty"`
	Permissions           string `json:"permissions,omitempty"`

	// Physio-specific fields
	LicenseNumber         string `json:"license_number,omitempty"`
	TreatmentAreas        string `json:"treatment_areas,omitempty"`
	Equipment             string `json:"equipment,omitempty"`
	Techniques            string `json:"techniques,omitempty"`
	InsuranceAccepted     string `json:"insurance_accepted,omitempty"`
	Affiliations          string `json:"affiliations,omitempty"`
}

// RoleProfileResponse represents the profile response for different roles
type RoleProfileResponse struct {
	UserID                uint    `json:"user_id"`
	Role                  string  `json:"role"`
	IsProfileComplete     bool    `json:"is_profile_complete"`
	CompletionPercentage  float64 `json:"completion_percentage"`
	
	// Role-specific data
	MemberProfile         *models.UserProfile    `json:"member_profile,omitempty"`
	TrainerProfile        *models.TrainerProfile `json:"trainer_profile,omitempty"`
	AdminProfile          *models.AdminProfile   `json:"admin_profile,omitempty"`
	PhysioProfile         *models.PhysioProfile  `json:"physio_profile,omitempty"`
}

// SetupRoleProfile sets up profile based on user role
func (s *RoleProfileService) SetupRoleProfile(userID uint, userRole string, req *RoleProfileSetupRequest) (*RoleProfileResponse, error) {
	// Update user phone number
	if err := s.updateUserPhone(userID, req.Phone); err != nil {
		return nil, err
	}

	var response *RoleProfileResponse
	var err error

	switch userRole {
	case "member":
		response, err = s.setupMemberProfile(userID, req)
	case "trainer":
		response, err = s.setupTrainerProfile(userID, req)
	case "admin":
		response, err = s.setupAdminProfile(userID, req)
	case "physio":
		response, err = s.setupPhysioProfile(userID, req)
	default:
		return nil, errors.New("unsupported user role")
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetRoleProfile retrieves profile based on user role
func (s *RoleProfileService) GetRoleProfile(userID uint, userRole string) (*RoleProfileResponse, error) {
	var response *RoleProfileResponse
	var err error

	switch userRole {
	case "member":
		response, err = s.getMemberProfile(userID)
	case "trainer":
		response, err = s.getTrainerProfile(userID)
	case "admin":
		response, err = s.getAdminProfile(userID)
	case "physio":
		response, err = s.getPhysioProfile(userID)
	default:
		return nil, errors.New("unsupported user role")
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

// CheckProfileCompletion checks profile completion based on role
func (s *RoleProfileService) CheckProfileCompletion(userID uint, userRole string) (*RoleProfileResponse, error) {
	profile, err := s.GetRoleProfile(userID, userRole)
	if err != nil {
		return nil, err
	}

	// Calculate completion percentage based on role
	completion := s.calculateRoleProfileCompletion(userRole, profile)
	profile.CompletionPercentage = completion
	profile.IsProfileComplete = completion >= 80.0

	return profile, nil
}

// Helper methods for different roles
func (s *RoleProfileService) setupMemberProfile(userID uint, req *RoleProfileSetupRequest) (*RoleProfileResponse, error) {
	var profile models.UserProfile
	s.db.Where("user_id = ?", userID).FirstOrCreate(&profile, models.UserProfile{UserID: userID})

	// Update profile fields
	profile.Height = req.Height
	profile.Weight = req.Weight
	profile.Age = req.Age
	profile.Gender = req.Gender
	profile.Goals = req.Goals
	profile.TargetWeight = req.TargetWeight
	profile.Timeline = req.Timeline
	profile.MedicalHistory = req.MedicalHistory
	profile.Allergies = req.Allergies
	profile.Medications = req.Medications
	profile.PhysioNeeds = req.PhysioNeeds
	profile.BodyFatPercentage = req.BodyFatPercentage
	profile.MuscleMass = req.MuscleMass
	profile.BodyMeasurements = req.BodyMeasurements
	profile.PreferredWorkoutTime = req.PreferredWorkoutTime
	profile.WorkoutDays = req.WorkoutDays
	profile.CommunicationPreference = req.CommunicationPreference

	// Calculate completion
	completion := s.calculateMemberProfileCompletion(&profile)
	profile.IsProfileComplete = completion >= 80.0

	if err := s.db.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "member",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		MemberProfile:        &profile,
	}, nil
}

func (s *RoleProfileService) setupTrainerProfile(userID uint, req *RoleProfileSetupRequest) (*RoleProfileResponse, error) {
	var profile models.TrainerProfile
	s.db.Where("user_id = ?", userID).FirstOrCreate(&profile, models.TrainerProfile{UserID: userID})

	// Update profile fields
	profile.Certifications = req.Certifications
	profile.Experience = req.Experience
	profile.Specializations = req.Specializations
	profile.Bio = req.Bio
	profile.Philosophy = req.Philosophy
	profile.Availability = req.Availability
	profile.SessionRates = req.SessionRates
	profile.PackageRates = req.PackageRates
	profile.Languages = req.Languages
	profile.Education = req.Education
	profile.Awards = req.Awards
	profile.PortfolioImages = req.PortfolioImages
	profile.BeforeAfterPhotos = req.BeforeAfterPhotos
	profile.Testimonials = req.Testimonials
	profile.PreferredContactMethod = req.PreferredContactMethod
	profile.ResponseTime = req.ResponseTime

	// Calculate completion
	completion := s.calculateTrainerProfileCompletion(&profile)
	profile.IsProfileComplete = completion >= 80.0

	if err := s.db.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "trainer",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		TrainerProfile:       &profile,
	}, nil
}

func (s *RoleProfileService) setupAdminProfile(userID uint, req *RoleProfileSetupRequest) (*RoleProfileResponse, error) {
	var profile models.AdminProfile
	s.db.Where("user_id = ?", userID).FirstOrCreate(&profile, models.AdminProfile{UserID: userID})

	// Update profile fields
	profile.AdminRole = req.AdminRole
	profile.AccessLevel = req.AccessLevel
	profile.Department = req.Department
	profile.EmergencyContact = req.EmergencyContact
	profile.EmergencyPhone = req.EmergencyPhone
	profile.OfficeLocation = req.OfficeLocation
	profile.Permissions = req.Permissions

	// Calculate completion
	completion := s.calculateAdminProfileCompletion(&profile)
	profile.IsProfileComplete = completion >= 80.0

	if err := s.db.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "admin",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		AdminProfile:         &profile,
	}, nil
}

func (s *RoleProfileService) setupPhysioProfile(userID uint, req *RoleProfileSetupRequest) (*RoleProfileResponse, error) {
	var profile models.PhysioProfile
	s.db.Where("user_id = ?", userID).FirstOrCreate(&profile, models.PhysioProfile{UserID: userID})

	// Update profile fields
	profile.LicenseNumber = req.LicenseNumber
	profile.Certifications = req.Certifications
	profile.Experience = req.Experience
	profile.Specializations = req.Specializations
	profile.TreatmentAreas = req.TreatmentAreas
	profile.Equipment = req.Equipment
	profile.Techniques = req.Techniques
	profile.Availability = req.Availability
	profile.SessionRates = req.SessionRates
	profile.InsuranceAccepted = req.InsuranceAccepted
	profile.Education = req.Education
	profile.Affiliations = req.Affiliations
	profile.Languages = req.Languages

	// Calculate completion
	completion := s.calculatePhysioProfileCompletion(&profile)
	profile.IsProfileComplete = completion >= 80.0

	if err := s.db.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "physio",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		PhysioProfile:        &profile,
	}, nil
}

// Get profile methods
func (s *RoleProfileService) getMemberProfile(userID uint) (*RoleProfileResponse, error) {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	completion := s.calculateMemberProfileCompletion(&profile)
	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "member",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		MemberProfile:        &profile,
	}, nil
}

func (s *RoleProfileService) getTrainerProfile(userID uint) (*RoleProfileResponse, error) {
	var profile models.TrainerProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	completion := s.calculateTrainerProfileCompletion(&profile)
	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "trainer",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		TrainerProfile:       &profile,
	}, nil
}

func (s *RoleProfileService) getAdminProfile(userID uint) (*RoleProfileResponse, error) {
	var profile models.AdminProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	completion := s.calculateAdminProfileCompletion(&profile)
	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "admin",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		AdminProfile:         &profile,
	}, nil
}

func (s *RoleProfileService) getPhysioProfile(userID uint) (*RoleProfileResponse, error) {
	var profile models.PhysioProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	completion := s.calculatePhysioProfileCompletion(&profile)
	return &RoleProfileResponse{
		UserID:               userID,
		Role:                 "physio",
		IsProfileComplete:    profile.IsProfileComplete,
		CompletionPercentage: completion,
		PhysioProfile:        &profile,
	}, nil
}

// Completion calculation methods
func (s *RoleProfileService) calculateRoleProfileCompletion(userRole string, profile *RoleProfileResponse) float64 {
	switch userRole {
	case "member":
		return s.calculateMemberProfileCompletion(profile.MemberProfile)
	case "trainer":
		return s.calculateTrainerProfileCompletion(profile.TrainerProfile)
	case "admin":
		return s.calculateAdminProfileCompletion(profile.AdminProfile)
	case "physio":
		return s.calculatePhysioProfileCompletion(profile.PhysioProfile)
	default:
		return 0.0
	}
}

func (s *RoleProfileService) calculateMemberProfileCompletion(profile *models.UserProfile) float64 {
	if profile == nil {
		return 0.0
	}

	fields := []bool{
		profile.Height > 0,
		profile.Weight > 0,
		profile.Age > 0,
		profile.Gender != "",
		profile.Goals != "",
		profile.TargetWeight > 0,
		profile.Timeline > 0,
		profile.MedicalHistory != "",
		profile.PreferredWorkoutTime != "",
		profile.WorkoutDays != "",
		profile.CommunicationPreference != "",
	}

	completed := 0
	for _, field := range fields {
		if field {
			completed++
		}
	}

	return float64(completed) / float64(len(fields)) * 100
}

func (s *RoleProfileService) calculateTrainerProfileCompletion(profile *models.TrainerProfile) float64 {
	if profile == nil {
		return 0.0
	}

	fields := []bool{
		profile.Certifications != "",
		profile.Experience > 0,
		profile.Specializations != "",
		profile.Bio != "",
		profile.Philosophy != "",
		profile.Availability != "",
		profile.SessionRates != "",
		profile.Languages != "",
		profile.Education != "",
		profile.PreferredContactMethod != "",
		profile.ResponseTime != "",
	}

	completed := 0
	for _, field := range fields {
		if field {
			completed++
		}
	}

	return float64(completed) / float64(len(fields)) * 100
}

func (s *RoleProfileService) calculateAdminProfileCompletion(profile *models.AdminProfile) float64 {
	if profile == nil {
		return 0.0
	}

	fields := []bool{
		profile.AdminRole != "",
		profile.AccessLevel != "",
		profile.Department != "",
		profile.EmergencyContact != "",
		profile.EmergencyPhone != "",
		profile.OfficeLocation != "",
		profile.Permissions != "",
	}

	completed := 0
	for _, field := range fields {
		if field {
			completed++
		}
	}

	return float64(completed) / float64(len(fields)) * 100
}

func (s *RoleProfileService) calculatePhysioProfileCompletion(profile *models.PhysioProfile) float64 {
	if profile == nil {
		return 0.0
	}

	fields := []bool{
		profile.LicenseNumber != "",
		profile.Certifications != "",
		profile.Experience > 0,
		profile.Specializations != "",
		profile.TreatmentAreas != "",
		profile.Availability != "",
		profile.SessionRates != "",
		profile.Education != "",
		profile.Affiliations != "",
		profile.Languages != "",
	}

	completed := 0
	for _, field := range fields {
		if field {
			completed++
		}
	}

	return float64(completed) / float64(len(fields)) * 100
}

// Helper method to update user phone
func (s *RoleProfileService) updateUserPhone(userID uint, phone string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("phone", phone).Error
}
