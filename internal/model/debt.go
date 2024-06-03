package model

import (
	"time"

	"gorm.io/datatypes"
)

// Debt represents debt model
type Debt struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	SessionID int64     `json:"session_id"`

	Name            string         `json:"name" gorm:"type:varchar(50)"`
	RemainingAmount float64        `json:"remaining_amount"`
	MonthlyPayment  float64        `json:"monthly_payment"`
	AnnualInterest  float64        `json:"annual_interest"`
	Type            string         `json:"type" gorm:"type:varchar(10);default:FIXED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
	PaymentDeadline datatypes.Date `json:"payment_deadline" gorm:"default:NULL"`

	ForecastPaidOffDate string `json:"forecast_paid_off_date" gorm:"type:varchar(50)"`

	Session *Session `json:"session,omitempty"`
}
