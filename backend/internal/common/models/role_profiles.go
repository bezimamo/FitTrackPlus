package models

import (
	"time"

	"gorm.io/gorm"
)

// TrainerProfile contains trainer-specific information
type TrainerProfile struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id" gorm:"uniqueIndex"`
	
	// Professional Information
	Certifications        string         `json:"certifications"` // JSON string of certifications
	Experience            int            `json:"experience"` // years of experience
	Specializations       string         `json:"specializations"` // JSON string of specializations
	Bio                   string         `json:"bio"` // About me section
	Philosophy            string         `json:"philosophy"` // Training philosophy
	
	// Availability & Pricing
	Availability          string         `json:"availability"` // JSON string of available hours
	SessionRates          string         `json:"session_rates"` // JSON string of pricing
	PackageRates          string         `json:"package_rates"` // JSON string of package pricing
	
	// Professional Details
	Languages             string         `json:"languages"` // JSON string of spoken languages
	Education             string         `json:"education"` // Educational background
	Awards                string         `json:"awards"` // Professional awards/recognition
	
	// Portfolio
	PortfolioImages       string         `json:"portfolio_images"` // JSON string of image URLs
	BeforeAfterPhotos     string         `json:"before_after_photos"` // JSON string of transformation photos
	Testimonials          string         `json:"testimonials"` // JSON string of client testimonials
	
	// Contact & Communication
	PreferredContactMethod string        `json:"preferred_contact_method"`
	ResponseTime          string         `json:"response_time"` // typical response time
	
	// Profile Completion
	IsProfileComplete     bool           `json:"is_profile_complete" gorm:"default:false"`
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// AdminProfile contains admin-specific information
type AdminProfile struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id" gorm:"uniqueIndex"`
	
	// Administrative Information
	AdminRole             string         `json:"admin_role"` // system_admin, gym_manager, content_manager
	AccessLevel           string         `json:"access_level"` // full_access, limited_access, read_only
	Department            string         `json:"department"` // IT, Management, Operations
	
	// Contact Information
	EmergencyContact      string         `json:"emergency_contact"`
	EmergencyPhone        string         `json:"emergency_phone"`
	OfficeLocation        string         `json:"office_location"`
	
	// Administrative Details
	Permissions           string         `json:"permissions"` // JSON string of granted permissions
	LastLogin             *time.Time     `json:"last_login"`
	LoginHistory          string         `json:"login_history"` // JSON string of recent logins
	
	// Profile Completion
	IsProfileComplete     bool           `json:"is_profile_complete" gorm:"default:false"`
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// PhysioProfile contains physiotherapist-specific information
type PhysioProfile struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id" gorm:"uniqueIndex"`
	
	// Professional Information
	LicenseNumber         string         `json:"license_number"`
	Certifications        string         `json:"certifications"` // JSON string of certifications
	Experience            int            `json:"experience"` // years of experience
	Specializations       string         `json:"specializations"` // JSON string of specializations
	
	// Medical Focus Areas
	TreatmentAreas        string         `json:"treatment_areas"` // JSON string of treatment areas
	Equipment             string         `json:"equipment"` // JSON string of available equipment
	Techniques            string         `json:"techniques"` // JSON string of treatment techniques
	
	// Availability & Pricing
	Availability          string         `json:"availability"` // JSON string of available hours
	SessionRates          string         `json:"session_rates"` // JSON string of pricing
	InsuranceAccepted     string         `json:"insurance_accepted"` // JSON string of accepted insurance
	
	// Professional Details
	Education             string         `json:"education"` // Educational background
	Affiliations          string         `json:"affiliations"` // Professional affiliations
	Languages             string         `json:"languages"` // JSON string of spoken languages
	
	// Profile Completion
	IsProfileComplete     bool           `json:"is_profile_complete" gorm:"default:false"`
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"-" gorm:"foreignKey:UserID"`
}
