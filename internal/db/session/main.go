package session

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
	"gorm.io/gorm"
)

// NewDB returns a new session database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Session{})}
}

// DB represents the client for session table
type DB struct {
	*dbutil.DB
}

// FindByRefreshToken queries for single session by refresh token
func (d *DB) FindByRefreshToken(db *gorm.DB, token string) (*model.Session, error) {
	rec := new(model.Session)
	if err := d.View(db, rec, "refresh_token = ?", token); err != nil {
		return nil, err
	}
	return rec, nil
}
