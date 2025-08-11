package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in our system
// GORM will automatically create a table named "users" for this struct
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'member'"` // member, trainer, physio, admin
	Phone     string         `json:"phone"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete

	// Relationships - these will be populated when we query the database
	Profile   *UserProfile   `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Plans     []UserPlan     `json:"plans,omitempty" gorm:"foreignKey:UserID"`
	Progress  []ProgressLog  `json:"progress,omitempty" gorm:"foreignKey:UserID"`
	Bookings  []Booking      `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
	Payments  []Payment      `json:"payments,omitempty" gorm:"foreignKey:UserID"`
}

// UserProfile contains additional user information
type UserProfile struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id" gorm:"uniqueIndex"`
	
	// Basic Information
	Height                float64        `json:"height"` // in cm
	Weight                float64        `json:"weight"` // in kg
	Age                   int            `json:"age"`
	Gender                string         `json:"gender"`
	
	// Fitness Goals
	Goals                 string         `json:"goals"` // JSON string of fitness goals
	TargetWeight          float64        `json:"target_weight"`
	Timeline              int            `json:"timeline"` // days to achieve goal
	
	// Medical Information
	MedicalHistory        string         `json:"medical_history"`
	Allergies             string         `json:"allergies"`
	Medications           string         `json:"medications"`
	PhysioNeeds           string         `json:"physio_needs"`
	
	// Physical Measurements
	BodyFatPercentage     float64        `json:"body_fat_percentage"`
	MuscleMass            float64        `json:"muscle_mass"`
	BodyMeasurements      string         `json:"body_measurements"` // JSON string of measurements
	
	// Profile Image
	ProfileImageURL       string         `json:"profile_image_url"`
	
	// Preferences
	PreferredWorkoutTime  string         `json:"preferred_workout_time"`
	WorkoutDays           string         `json:"workout_days"` // JSON string of days
	CommunicationPreference string       `json:"communication_preference"`
	
	// Profile Completion
	IsProfileComplete     bool           `json:"is_profile_complete" gorm:"default:false"`
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// Plan represents a fitness plan template
type Plan struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	GoalType    string         `json:"goal_type"` // lose_weight, gain_muscle, flexibility, rehab
	PlanType    string         `json:"plan_type"` // fitness, diet, physio
	Exercises   string         `json:"exercises"` // JSON string of exercises
	Diet        string         `json:"diet"`      // JSON string of diet plan
	PhysioExercises string     `json:"physio_exercises"` // JSON string of physio exercises
	Duration    int            `json:"duration"` // in days
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	UserPlans []UserPlan `json:"user_plans,omitempty" gorm:"foreignKey:PlanID"`
}

// UserPlan links users to their assigned plans
type UserPlan struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	PlanID     uint           `json:"plan_id"`
	Status     string         `json:"status" gorm:"default:'active'"` // active, completed, paused
	AssignedAt time.Time      `json:"assigned_at"`
	CompletedAt *time.Time    `json:"completed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Plan Plan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
}

// ProgressLog tracks user progress over time
type ProgressLog struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id"`
	Weight          float64        `json:"weight"`
	Measurements    string         `json:"measurements"` // JSON string of body measurements
	WorkoutCompletion string       `json:"workout_completion"` // JSON string of completed workouts
	PhysioProgress  string         `json:"physio_progress"` // JSON string of physio progress
	Notes           string         `json:"notes"`
	LoggedAt        time.Time      `json:"logged_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Booking represents a training or physiotherapy session
type Booking struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	TrainerID   *uint          `json:"trainer_id"` // Can be null for physio sessions
	PhysioID    *uint          `json:"physio_id"`  // Can be null for training sessions
	SessionDate time.Time      `json:"session_date"`
	SessionType string         `json:"session_type"` // training, physio
	Status      string         `json:"status" gorm:"default:'pending'"` // pending, approved, completed, cancelled
	Notes       string         `json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Trainer *User `json:"trainer,omitempty" gorm:"foreignKey:TrainerID"`
	Physio  *User `json:"physio,omitempty" gorm:"foreignKey:PhysioID"`
}

// Payment tracks user payments
type Payment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	Amount     float64        `json:"amount"`
	Currency   string         `json:"currency" gorm:"default:'ETB'"`
	ChapaRef   string         `json:"chapa_ref"`
	Status     string         `json:"status" gorm:"default:'pending'"` // pending, completed, failed
	PaymentDate *time.Time    `json:"payment_date"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
} 