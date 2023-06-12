package model

import (
	"time"

	"gorm.io/gorm"
)

// Session represents the session model
// swagger:model
type Session struct {
	Base
	Code      string `json:"code" gorm:"type:varchar(20);unique_index"`
	IPAddress string `json:"ip_address" gorm:"type:varchar(45)" `
	UserAgent string `json:"user_agent" gorm:"type:text"`

	RefreshToken string     `json:"-" gorm:"type:varchar(100);unique_index"`
	LastLogin    *time.Time `json:"last_login"`

	TotalAllIncome           float64 `json:"total_all_income"`
	TotalAllExpense          float64 `json:"total_all_expense"`
	TotalMonthlyPaymentDebt  float64 `json:"total_monthly_payment_debt"`
	TotalEssentialExpense    float64 `json:"total_essential_expense"`
	TotalNonEssentialExpense float64 `json:"total_non_essential_expense"`
	MonthlyNetFlow           float64 `json:"monthly_net_flow"` // important

	CurrentBalance float64 `json:"current_balance"`

	ActualEmergencyFund   float64 `json:"actual_emergency_fund"`
	ExpectedEmergencyFund float64 `json:"expected_emergency_fund"`
	ActualRainydayFund    float64 `json:"actual_rainyday_fund"`
	ExpectedRainydayFund  float64 `json:"expected_rainyday_fund"`
	FunFund               float64 `json:"fun_fund"`
	Investment            float64 `json:"investment"`
	RetirementPlan        float64 `json:"retirement_plan"`

	IsAchivedEmergencyFund  bool `json:"is_achived_emergency_fund"`
	IsAchivedRainydayFund   bool `json:"is_achived_rainyday_fund"`
	IsAchivedInvestment     bool `json:"is_achived_investment"`
	IsAchivedRetirementPlan bool `json:"is_achived_retirement_plan"`

	Status      string `json:"status" gorm:"type:varchar(10)"`
	Description string `json:"description" gorm:"-"`

	Incomes  []*Income  `json:"incomes,omitempty"`
	Expenses []*Expense `json:"expenses,omitempty"`
	Debts    []*Debt    `json:"debts,omitempty"`

	DataSets []*DataSet `json:"data_sets,omitempty" gorm:"-"`
}

// DataSet represents the dataset model
// swagger:model
type DataSet struct {
	Group string  `json:"group"`
	Key   string  `json:"key"`
	Asset float64 `json:"asset"`
	Debt  float64 `json:"debt"`
}

// "total_all_income":           totalIncome,
// "total_all_expense":          totalExpense,
// "total_monthly_payment_debt": totalMonthlyPaymentDebt,
// "monthly_net_flow":           monthlyNetFlow,
// "status":                     status,
// "expected_emergency_fund":    expectedEmergencyFund,
// "expected_rainyday_fund":     expectedRainydayFund,
// "actual_emergency_fund":      actualEmergencyFund,
// "actual_rainyday_fund":       actualRainydayFund,
// "fun_fund":                   funFund,
// "retirement_plan":            retirementPlan,
// "is_achived_emergency_fund":  isAchivedEmergencyFund,
// "is_achived_rainyday_fund":   isAchivedRainydayFund,
// "is_achived_investment":      isAchivedInvestment,
// "is_achived_retirement_plan": isAchivedRetirementPlan,

// DataNode represents the data each node
// swagger:model
type DataNode struct {
	SessionID int64  `json:"session_id"`
	NodeName  string `json:"node_name"`

	CurrentAsset            float64 `json:"current_asset"`
	TotalAllIncome          float64 `json:"total_all_income"`
	TotalAllExpense         float64 `json:"total_all_expense"`
	TotalMonthlyPaymentDebt float64 `json:"total_monthly_payment_debt"`
	MonthlyNetFlow          float64 `json:"monthly_net_flow"`
	Status                  string  `json:"status"`
	Descrtiption            string  `json:"description"`
	ExpectedEmergencyFund   float64 `json:"expected_emergency_fund"`
	ExpectedRainydayFund    float64 `json:"expected_rainyday_fund"`
	ActualEmergencyFund     float64 `json:"actual_emergency_fund"`
	ActualRainydayFund      float64 `json:"actual_rainyday_fund"`
	FunFund                 float64 `json:"fun_fund"`
	RetirementPlan          float64 `json:"retirement_plan"`
	IsAchivedEmergencyFund  bool    `json:"is_achived_emergency_fund"`
	IsAchivedRainydayFund   bool    `json:"is_achived_rainyday_fund"`
	IsAchivedInvestment     bool    `json:"is_achived_investment"`
	IsAchivedRetirementPlan bool    `json:"is_achived_retirement_plan"`
}

// DataDebtNode represents the data of debt for each node
// swagger:model
type DataDebtNode struct {
	SessionID int64  `json:"session_id"`
	NodeName  string `json:"node_name"`
	DebtID    int64  `json:"debt_id"`
	Index     int    `json:"index"`

	RemainingAmount float64 `json:"remaining_amount"`
	MonthlyPayment  float64 `json:"monthly_payment"`

	IsEligiblePaidOff bool `json:"is_eligible_paid_off"`
	IsPaidOff         bool `json:"is_paid_off"`
}

// AfterSave to run after save
func (s *Session) AfterSave(tx *gorm.DB) (err error) {
	// do something here
	return
}

// AfterUpdate to run after update
func (s *Session) AfterUpdate(tx *gorm.DB) (err error) {
	// do something here
	return
}

// Custom status
const (
	SessionStatusDefault = "DEFAULT"
	SessionStatusBD      = "BD"    // Budget Deficit
	SessionStatusPC2PC   = "PC2PC" // Pay Check to Pay Check
	SessionStatusLFF     = "LFF"   // Limited financial flexibility
	SessionStatusGFF     = "GFF"   // Good financial flexibility

	SessionStatusDescriptionDefault = "Our BA has not analyzed your financial situation yet. Please wait for a decade."
	SessionStatusDescriptionBD      = "You are spending more money than they are earning. This is not sustainable in the long term and can lead to financial problems. It is important to take steps to increase your income or decrease your expenses in order to bring your budget back into balance."
	SessionStatusDescriptionPC2PC   = "You are earning just enough money to cover your expenses, but you do not have any extra money left over at the end of the month. This can be a stressful financial situation since you have no financial cushion in case of an emergency or unexpected expense. It is important to find ways to increase income or decrease expenses in order to break out of this cycle and build up savings."
	SessionStatusDescriptionLFF     = "You have some extra money left over at the end of the month, but not enough to save or invest. This can make it difficult for you to respond to unexpected expenses or changes in income."
	SessionStatusDescriptionGFF     = "Your monthly net income exceeds your essential expenses. This provides you with extra money at the end of each month that you can save or invest. This extra financial cushion can help you respond to unexpected expenses or changes in income. Consider saving toward the Emergency and Rainy Day Fund if you haven’t done so."

	DataSetTypeAsset = "asset"
	DataSetTypeDebt  = "debt"
)

// Session status
var (
	SessionStatusDescriptions = map[string]string{
		SessionStatusDefault: SessionStatusDescriptionDefault,
		SessionStatusBD:      SessionStatusDescriptionBD,
		SessionStatusPC2PC:   SessionStatusDescriptionPC2PC,
		SessionStatusLFF:     SessionStatusDescriptionLFF,
		SessionStatusGFF:     SessionStatusDescriptionGFF,
	}
)
