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

	TotalIncome              float64 `json:"total_income"`
	TotalEssentialExpense    float64 `json:"total_essential_expense"`
	TotalNonEssentialExpense float64 `json:"total_non_essential_expense"`
	MonthlyPaymentDebt       float64 `json:"monthly_payment_debt"`
	MonthlyNetFlow           float64 `json:"monthly_net_flow"`

	EmergencyAchieved float64 `json:"emergency_achieved"`
	RainyFundAchieved float64 `json:"rainy_fund_achieved"`
	FunFund           float64 `json:"fun_fund"`
	Investment        float64 `json:"investment"`
	CurrentRetirement float64 `json:"current_retirement"`

	Status string `json:"status" gorm:"type:varchar(10)"`

	Incomes  []*Income  `json:"incomes,omitempty"`
	Expenses []*Expense `json:"expenses,omitempty"`
	Debts    []*Debt    `json:"debts,omitempty"`
}

// Custom status
const (
	SessionStatusBD    = "BD"    // Budget Deficit
	SessionStatusPC2PC = "PC2PC" // Pay Check to Pay Check
	SessionStatusLFF   = "LFF"   // Limited financial flexibility
	SessionStatusGFF   = "GFF"   // Good financial flexibility

	SessionStatusDescriptionBD    = "You are spending more money than they are earning. This is not sustainable in the long term and can lead to financial problems. It is important to take steps to increase your income or decrease your expenses in order to bring your budget back into balance."
	SessionStatusDescriptionPC2PC = "You are earning just enough money to cover your expenses, but you do not have any extra money left over at the end of the month. This can be a stressful financial situation since you have no financial cushion in case of an emergency or unexpected expense. It is important to find ways to increase income or decrease expenses in order to break out of this cycle and build up savings."
	SessionStatusDescriptionLFF   = "You have some extra money left over at the end of the month, but not enough to save or invest. This can make it difficult for you to respond to unexpected expenses or changes in income."
	SessionStatusDescriptionGFF   = "Your monthly net income exceeds your essential expenses. This provides you with extra money at the end of each month that you can save or invest. This extra financial cushion can help you respond to unexpected expenses or changes in income. Consider saving toward the Emergency and Rainy Day Fund if you haven’t done so."
)

// Session status
var (
	SessionStatusDescriptions = map[string]string{
		SessionStatusBD:    SessionStatusDescriptionBD,
		SessionStatusPC2PC: SessionStatusDescriptionPC2PC,
		SessionStatusLFF:   SessionStatusDescriptionLFF,
		SessionStatusGFF:   SessionStatusDescriptionGFF,
	}
)
