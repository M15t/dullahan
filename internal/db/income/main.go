package income

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
	"gorm.io/gorm"
)

// NewDB returns a new income database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Income{})}
}

// DB represents the client for income table
type DB struct {
	*dbutil.DB
}

// SumTotalIncome get sum total income
func (d *DB) SumTotalIncome(db *gorm.DB, totalIncome *float64, sessionID int64) error {
	return db.Raw(`SELECT SUM(amount) FROM incomes WHERE session_id = ?`, sessionID).Scan(totalIncome).Error
}
