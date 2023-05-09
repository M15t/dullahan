package model

import "time"

// Income represents income model
// swagger:model
type Income struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	SessionID int64   `json:"-"`
	Amount    float64 `json:"amount"`
	Name      string  `json:"name" gorm:"type:varchar(100)"`
	Type      string  `json:"type" gorm:"type:varchar(10);default:MONTHLY"` // MONTHLY, PASSIVE

	Session *Session `json:"session,omitempty"`
}
