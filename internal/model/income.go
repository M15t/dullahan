package model

// Income represents income model
// swagger:model
type Income struct {
	Base
	SessionID int64   `json:"session_id"`
	Amount    float64 `json:"amount"`
	Name      string  `json:"name" gorm:"type:varchar(100)"`
	Type      string  `json:"type" gorm:"type:varchar(10);default:MONTHLY"` // MONTHLY, PASSIVE

	Session *Session `json:"session,omitempty"`
}
