package model

// User represents the user model
// swagger:model
type User struct {
	Base
	Name  string `gorm:"type:varchar(100);not null" json:"name"`
	Email string `gorm:"type:varchar(150);not null;unique" json:"email"`
	Role  string `gorm:"type:varchar(10);not null" json:"role"`
}
