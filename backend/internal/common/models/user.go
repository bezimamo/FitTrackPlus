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

// PlanRequest represents a member's request for a plan (pending approval)
type PlanRequest struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"` // Member requesting the plan
	PlanID      uint           `json:"plan_id"` // Plan being requested
	Reason      string         `json:"reason"`  // Why they want this plan
	Status      string         `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
	RequestedAt time.Time      `json:"requested_at"`
	ReviewedAt  *time.Time     `json:"reviewed_at,omitempty"`
	ReviewedBy  *uint          `json:"reviewed_by,omitempty"` // Admin who reviewed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User       User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Plan       Plan  `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	ReviewedByUser *User `json:"reviewed_by_user,omitempty" gorm:"foreignKey:ReviewedBy"`
}

// TrainerAssignment represents the assignment of a trainer to manage a member's plan
type TrainerAssignment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	TrainerID  uint           `json:"trainer_id"` // Trainer assigned to manage
	MemberID   uint           `json:"member_id"`  // Member receiving the plan
	PlanID     uint           `json:"plan_id"`    // Plan being managed
	UserPlanID uint           `json:"user_plan_id"` // Reference to the UserPlan
	Status     string         `json:"status" gorm:"default:'active'"` // active, completed, paused
	AssignedAt time.Time      `json:"assigned_at"`
	AssignedBy uint           `json:"assigned_by"` // Admin who made assignment
	Notes      string         `json:"notes"`       // Assignment notes
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Trainer    User     `json:"trainer,omitempty" gorm:"foreignKey:TrainerID"`
	Member     User     `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Plan       Plan     `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	UserPlan   UserPlan `json:"user_plan,omitempty" gorm:"foreignKey:UserPlanID"`
	AssignedByUser User `json:"assigned_by_user,omitempty" gorm:"foreignKey:AssignedBy"`
}

// Exercise represents individual fitness exercises
type Exercise struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"` // strength, cardio, flexibility, balance
	MuscleGroup string         `json:"muscle_group"` // chest, back, legs, shoulders, arms, core, full_body
	Difficulty  string         `json:"difficulty"` // beginner, intermediate, advanced
	Equipment   string         `json:"equipment"` // bodyweight, dumbbells, barbell, machine, etc.
	VideoURL    string         `json:"video_url"`
	ImageURL    string         `json:"image_url"`
	Instructions string        `json:"instructions"`
	Tips        string         `json:"tips"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	WorkoutExercises []WorkoutExercise `json:"workout_exercises,omitempty" gorm:"foreignKey:ExerciseID"`
}

// Workout represents a complete workout session
type Workout struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	MemberID    uint           `json:"member_id"` // Who this workout is for
	TrainerID   uint           `json:"trainer_id"` // Who created this workout
	PlanID      *uint          `json:"plan_id,omitempty"` // Optional: linked to a plan
	Difficulty  string         `json:"difficulty"` // beginner, intermediate, advanced
	EstimatedDuration int      `json:"estimated_duration"` // minutes
	Notes       string         `json:"notes"` // Instructions for the member
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Member      User            `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Trainer     User            `json:"trainer,omitempty" gorm:"foreignKey:TrainerID"`
	Plan        *Plan           `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Exercises   []WorkoutExercise `json:"exercises,omitempty" gorm:"foreignKey:WorkoutID"`
	WorkoutLogs []WorkoutLog    `json:"workout_logs,omitempty" gorm:"foreignKey:WorkoutID"`
}

// WorkoutExercise links exercises to workouts with specific parameters
type WorkoutExercise struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	WorkoutID   uint           `json:"workout_id"`
	ExerciseID  uint           `json:"exercise_id"`
	Order       int            `json:"order"` // Exercise order in workout
	Sets        int            `json:"sets"` // Number of sets
	Reps        int            `json:"reps"` // Reps per set
	Weight      *float64       `json:"weight,omitempty"` // Weight in kg (optional)
	Duration    *int           `json:"duration,omitempty"` // Duration in seconds (for cardio)
	RestTime    int            `json:"rest_time"` // Rest time between sets in seconds
	Notes       string         `json:"notes"` // Specific notes for this exercise
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relationships
	Workout     Workout        `json:"workout,omitempty" gorm:"foreignKey:WorkoutID"`
	Exercise    Exercise       `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
}

// WorkoutLog tracks when members complete workouts
type WorkoutLog struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id"`
	WorkoutID       uint           `json:"workout_id"`
	CompletedAt     time.Time      `json:"completed_at"`
	ActualDuration  int            `json:"actual_duration"` // minutes
	ExercisesCompleted int         `json:"exercises_completed"`
	TotalSets       int            `json:"total_sets"`
	Rating          int            `json:"rating"` // 1-5 stars
	Notes           string         `json:"notes"` // How it felt, any issues
	CaloriesBurned  *int           `json:"calories_burned,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relationships
	User           User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Workout        Workout        `json:"workout,omitempty" gorm:"foreignKey:WorkoutID"`
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