package expense

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
	"gorm.io/gorm"
)

// NewDB returns a new expense database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Expense{})}
}

// DB represents the client for expense table
type DB struct {
	*dbutil.DB
}

// SumExpenseByType get sum expense by type
func (d *DB) SumExpenseByType(db *gorm.DB, totaExpense *float64, dataType string, sessionID int64) error {
	return db.Raw(`SELECT SUM(amount) FROM expenses WHERE session_id = ? AND type = ?`, sessionID, dataType).Scan(totaExpense).Error
}
