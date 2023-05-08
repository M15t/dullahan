package auth

import (
	"time"

	"dullahan/config"
	"dullahan/internal/db"
)

// New creates new auth service
func New(db *db.Service, jwt JWT, cr Crypter, cfg *config.Configuration) *Auth {
	return &Auth{
		db:  db,
		jwt: jwt,
		cr:  cr,
		cfg: cfg,
	}
}

// Auth represents auth application service
type Auth struct {
	db  *db.Service
	jwt JWT
	cr  Crypter
	cfg *config.Configuration
}

// JWT represents token generator (jwt) interface
type JWT interface {
	GenerateToken(map[string]interface{}, *time.Time) (string, int, error)
}

// Crypter represents security interface
type Crypter interface {
	HashPassword(password string) string
	CompareHashAndPassword(string, string) bool
	UID() string
	NanoID() (string, error)
}
