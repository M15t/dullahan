package db

import (
	debtDB "dullahan/internal/db/debt"
	expenseDB "dullahan/internal/db/expense"
	incomeDB "dullahan/internal/db/income"
	sessionDB "dullahan/internal/db/session"

	"gorm.io/gorm"
)

// Service provides all databases
type Service struct {
	GDB     *gorm.DB
	Session *sessionDB.DB
	Income  *incomeDB.DB
	Expense *expenseDB.DB
	Debt    *debtDB.DB
}

// New creates db service
func New(db *gorm.DB) *Service {
	return &Service{
		GDB:     db,
		Session: sessionDB.NewDB(),
		Income:  incomeDB.NewDB(),
		Expense: expenseDB.NewDB(),
		Debt:    debtDB.NewDB(),
	}
}
