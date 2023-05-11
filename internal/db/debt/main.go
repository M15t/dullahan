package debt

import (
	"dullahan/internal/model"

	dbutil "github.com/M15t/ghoul/pkg/util/db"
	"gorm.io/gorm"
)

// NewDB returns a new debt database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(&model.Debt{})}
}

// DB represents the client for debt table
type DB struct {
	*dbutil.DB
}

// SumTotalMonthlyPaymentDebt get sum total monthly payment
func (d *DB) SumTotalMonthlyPaymentDebt(db *gorm.DB, totalMonthlyPayment *float64, sessionID int64) error {
	return db.Raw(`SELECT SUM(monthly_payment) FROM debts WHERE session_id = ?`, sessionID).Scan(totalMonthlyPayment).Error
}
