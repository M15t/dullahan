package model

import "time"

// Session represents the session model
// swagger:model
type Session struct {
	Base
	Code      string `json:"code" gorm:"type:varchar(20);unique_index"`
	IPAddress string `json:"ip_address" gorm:"type:varchar(45)" `
	UserAgent string `json:"user_agent" gorm:"type:text"`

	RefreshToken string     `json:"-" gorm:"type:varchar(100);unique_index"`
	LastLogin    *time.Time `json:"last_login"`
}
