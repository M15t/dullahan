package expense

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
)

// NewDB returns a new expense database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Expense{})}
}

// DB represents the client for expense table
type DB struct {
	*dbutil.DB
}
