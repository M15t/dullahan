package debt

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
)

// NewDB returns a new debt database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Debt{})}
}

// DB represents the client for debt table
type DB struct {
	*dbutil.DB
}
