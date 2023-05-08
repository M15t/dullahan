package income

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
)

// NewDB returns a new income database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Income{})}
}

// DB represents the client for income table
type DB struct {
	*dbutil.DB
}
