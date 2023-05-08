package model

import (
	"time"
)

// Base contains common fields for all models
// Do not use gorm.Model because of uint ID
type Base struct {
	// The primary key of the record
	ID int64 `json:"id" gorm:"primary_key"`
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
}

// BaseWithoutID contains common fields for all models that without ID included
// Do not use gorm.Model because of uint ID
type BaseWithoutID struct {
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
}
