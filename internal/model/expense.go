package model

// Expense represents expense model
// swagger:model
type Expense struct {
	Base
	SessionID int64   `json:"session_id"`
	Amount    float64 `json:"amount"`
	Name      string  `json:"name" gorm:"type:varchar(100)"`
	Type      string  `json:"type" gorm:"type:varchar(15);default:ESSENTIAL"` // ESSENTIAL, NON_ESSENTIAL

	Session *Session `json:"session,omitempty"`
}
