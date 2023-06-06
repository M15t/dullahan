package model

import "time"

// Expense represents expense model
// swagger:model
type Expense struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	SessionID int64     `json:"session_id"`

	Amount float64 `json:"amount"`
	Name   string  `json:"name" gorm:"type:varchar(100)"`
	Type   string  `json:"type" gorm:"type:varchar(15);default:ESSENTIAL"` // ESSENTIAL, NON_ESSENTIAL

	Session *Session `json:"session,omitempty"`
}

// Custom const
const (
	ExpenseTypeEssential    = "ESSENTIAL"
	ExpenseTypeNonEssential = "NON_ESSENTIAL"
)
