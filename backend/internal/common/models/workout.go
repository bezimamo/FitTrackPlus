package models

import (
	"time"
)

// Workout represents a workout template or session
type Workout struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Category    string    `json:"category" gorm:"not null"` // strength, cardio, hiit, flexibility, etc.
	Difficulty  string    `json:"difficulty" gorm:"not null"` // beginner, intermediate, advanced
	Duration    int       `json:"duration"` // in minutes
	Equipment   string    `json:"equipment"` // required equipment
	IsTemplate  bool      `json:"is_template" gorm:"default:true"` // true for templates, false for actual sessions
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"` // user ID who created it
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	CreatedByUser *User            `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBy"`
	Exercises     []WorkoutExercise `json:"exercises,omitempty" gorm:"foreignKey:WorkoutID"`
	UserWorkouts  []UserWorkout    `json:"user_workouts,omitempty" gorm:"foreignKey:WorkoutID"`
}

// WorkoutExercise represents an exercise within a workout
type WorkoutExercise struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	WorkoutID    uint    `json:"workout_id" gorm:"not null"`
	ExerciseID   uint    `json:"exercise_id" gorm:"not null"`
	Order        int     `json:"order" gorm:"not null"` // order of exercise in workout
	Sets         int     `json:"sets" gorm:"not null"`
	Reps         int     `json:"reps"` // can be 0 for time-based exercises
	Weight       float64 `json:"weight"` // in kg, can be 0 for bodyweight
	Duration     int     `json:"duration"` // in seconds, for time-based exercises
	RestTime     int     `json:"rest_time"` // rest between sets in seconds
	Notes        string  `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	Workout  *Workout  `json:"workout,omitempty" gorm:"foreignKey:WorkoutID"`
	Exercise *Exercise `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
}

// UserWorkout represents a user's workout session
type UserWorkout struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	WorkoutID    uint      `json:"workout_id" gorm:"not null"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Status       string    `json:"status" gorm:"default:'in_progress'"` // in_progress, completed, paused, cancelled
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	User         *User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Workout      *Workout             `json:"workout,omitempty" gorm:"foreignKey:WorkoutID"`
	ExerciseSets []UserWorkoutExercise `json:"exercise_sets,omitempty" gorm:"foreignKey:UserWorkoutID"`
}

// UserWorkoutExercise represents a user's performance for a specific exercise in a workout
type UserWorkoutExercise struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserWorkoutID uint      `json:"user_workout_id" gorm:"not null"`
	ExerciseID    uint      `json:"exercise_id" gorm:"not null"`
	SetNumber     int       `json:"set_number" gorm:"not null"`
	Reps          int       `json:"reps"`
	Weight        float64   `json:"weight"`
	Duration      int       `json:"duration"` // in seconds
	RestTime      int       `json:"rest_time"` // actual rest taken
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relationships
	UserWorkout *UserWorkout `json:"user_workout,omitempty" gorm:"foreignKey:UserWorkoutID"`
	Exercise    *Exercise    `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
}

// WorkoutCategory represents workout categories
type WorkoutCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkoutDifficulty represents workout difficulty levels
type WorkoutDifficulty struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
