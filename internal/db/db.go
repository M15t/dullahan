package db

import (
	sessionDB "dullahan/internal/db/session"

	"gorm.io/gorm"
)

// Service provides all databases
type Service struct {
	GDB     *gorm.DB
	Session *sessionDB.DB
}

// New creates db service
func New(db *gorm.DB) *Service {
	return &Service{
		GDB:     db,
		Session: sessionDB.NewDB(),
	}
}
