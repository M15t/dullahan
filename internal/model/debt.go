package model

import "gorm.io/datatypes"

// Debt represents debt model
type Debt struct {
	Base
	SessionID       int64          `json:"session_id"`
	Name            string         `json:"name" gorm:"type:varchar(50)"`
	RemainingAmount float64        `json:"remaining_amount"`
	MonthlyPayment  float64        `json:"monthly_payment"`
	AnnualInterest  float64        `json:"annual_interest"`
	Type            string         `json:"type" gorm:"type:varchar(10);default:FIXED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
	PaymentDeadline datatypes.Date `json:"payment_deadline"`

	DebtPaidOffEachMonth     float64 `json:"debt_paid_off_each_month"`
	InterestPaidOffEachMonth float64 `json:"interest_paid_off_each_month"`

	Session *Session `json:"session,omitempty"`
}
