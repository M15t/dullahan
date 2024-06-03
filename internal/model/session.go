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

	ActualRainydayFund   float64 `json:"actual_rainyday_fund"`
	ExpectedRainydayFund float64 `json:"expected_rainyday_fund"`

	ActualFunFund   float64 `json:"actual_fun_fund"`
	ExpectedFunFund float64 `json:"expected_fun_fund"`

	Investment     float64 `json:"investment"`
	RetirementPlan float64 `json:"retirement_plan"`

	IsAchivedEmergencyFund  bool `json:"is_achived_emergency_fund"`
	IsAchivedRainydayFund   bool `json:"is_achived_rainyday_fund"`
	IsAchivedInvestment     bool `json:"is_achived_investment"`
	IsAchivedRetirementPlan bool `json:"is_achived_retirement_plan"`

	ForecastEmergencyBudgetFilledDate string `json:"forecast_emergency_budget_filled_date" gorm:"type:varchar(50)"`
	ForecastRainydayBudgetFilledDate  string `json:"forecast_rainyday_budget_filled_date" gorm:"type:varchar(50)"`
	ForecastStartInvestingDate        string `json:"forecast_start_investing_date" gorm:"type:varchar(50)"`
	ForecastFinancialFreedomDate      string `json:"forecast_financial_freedom_date" gorm:"type:varchar(50)"`
	ForecastMillionaireDate           string `json:"forecast_millionaire_date" gorm:"type:varchar(50)"`
	ForecastBankrupt                  string `json:"forecast_bankrupt" gorm:"type:varchar(50)"`

	Status      string `json:"status" gorm:"type:varchar(10)"`
	FullStatus  string `json:"full_status" gorm:"-"`
	Description string `json:"description" gorm:"-"`

	NextNYears int `json:"next_n_years" gorm:"-"`

	Incomes  []*Income  `json:"incomes,omitempty"`
	Expenses []*Expense `json:"expenses,omitempty"`
	Debts    []*Debt    `json:"debts,omitempty"`

	// DataLinecharts []*LineChart `json:"data_linecharts,omitempty" gorm:"-"`
	// DataTimelines  []*Timeline  `json:"data_timelines,omitempty" gorm:"-"`
}

// LineChart represents the dataset model
// swagger:model
type LineChart struct {
	Group string  `json:"group"`
	Key   string  `json:"key"`
	Asset float64 `json:"asset"`
	Debt  float64 `json:"debt"`
}

// Timeline represents the data timeline
// swagger:model
type Timeline struct {
	Event       string    `json:"event"`
	Date        string    `json:"date"`
	Datetime    time.Time `json:"-"`
	Description string    `json:"description"`
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

	CurrentAsset              float64 `json:"current_asset"`
	TotalAllIncome            float64 `json:"total_all_income"`
	TotalAllExpense           float64 `json:"total_all_expense"`
	TotalMonthlyPaymentDebt   float64 `json:"total_monthly_payment_debt"`
	MonthlyNetFlow            float64 `json:"monthly_net_flow"`
	Status                    string  `json:"status"`
	Descrtiption              string  `json:"description"`
	ExpectedEmergencyFund     float64 `json:"expected_emergency_fund"`
	ExpectedRainydayFund      float64 `json:"expected_rainyday_fund"`
	ActualEmergencyFund       float64 `json:"actual_emergency_fund"`
	ActualRainydayFund        float64 `json:"actual_rainyday_fund"`
	ActualFunFund             float64 `json:"actual_fun_fund"`
	ExpectFunFund             float64 `json:"expected_fun_fund"`
	RetirementPlan            float64 `json:"retirement_plan"`
	IsAchivedEmergencyFund    bool    `json:"is_achived_emergency_fund"`
	IsAchivedRainydayFund     bool    `json:"is_achived_rainyday_fund"`
	IsAchivedInvestment       bool    `json:"is_achived_investment"`
	IsAchivedRetirementPlan   bool    `json:"is_achived_retirement_plan"`
	IsAchivedFinancialFreedom bool    `json:"is_achived_financial_freedom"`
	IsPaidAllDebt             bool    `json:"is_paid_all_debt"`
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
	SessionStatusDescriptionBD      = "You are spending more money than you earn. This is not sustainable in the long term and can lead to financial problems. It is important to take steps to increase your income or decrease your expenses in order to bring your budget back into balance."
	SessionStatusDescriptionPC2PC   = "You are earning just enough money to cover your expenses, but you do not have any extra money left over at the end of the month. This can be a stressful financial situation since you have no financial cushion in case of an emergency or unexpected expense. It is important to find ways to increase income or decrease expenses in order to break out of this cycle and build up savings."
	SessionStatusDescriptionLFF     = "You have some extra money left over at the end of the month, but not enough to save or invest. This can make it difficult for you to respond to unexpected expenses or changes in income."
	SessionStatusDescriptionGFF     = "Your monthly net income exceeds your essential expenses. This provides you with extra money at the end of each month that you can save or invest. This extra financial cushion can help you respond to unexpected expenses or changes in income. Consider saving toward the Emergency and Rainy Day Fund if you haven't done so."

	DatasetTypeAsset = "asset"
	DatasetTypeDebt  = "debt"

	SessionTitleForecastEmergencyBudgetFilled = "Emergency Budget Filled"
	SessionTitleForecastRainydayBudgetFilled  = "Rainy Day Budget Filled"
	SessionTitleForecastStartInvesting        = "Start Investing"
	SessionTitleForecastFinancialFreedom      = "Financial Freedom"
	SessionTitleForecastBankrupt              = "Watch out, your pocket is empty!!!"
	SessionTitleForecastMillionaire           = "Millionaire!!!"

	SessionDescriptionForecastEmergencyBudgetFilled = "You just achieved your Emergency Budget, now you will be feeling at ease in case any emergency happened"
	SessionDescriptionForecastRainydayBudgetFilled  = "Your Rainy Day budget achieved, well done, your Finance Journey will start getting easier from this point"
	SessionDescriptionForecastStartInvesting        = "You can now start Investing, whether by deposit to the bank, buying ETF funds. You should start doing research and let the money work for you. Averagely, investing in safe option can give you around 10%-11% annual interest rate. From this point forward, we will accumulate your asset as if you are investing to let you see the power of compound interest. However, please make sure you have prepared your knowledge in investing fields first before considering it into action."
	SessionDescriptionForecastFinancialFreedom      = "You are now Financially Free, you can now do almost whatever you want, maybe finding a job you really like, travel the world, or enjoy life a little. This doesn't mean the end of your financial journey though, it’s just improve your life quality from here by giving you more options to choose. Never stop trying to keep a good financial performance."
	SessionDescriptionForecastBankrupt              = "Oops, looks like you are running on a budget deficit which causes your budget to reach. We get it, life is tough, you have bills to pay, people to take care and not to mention time for yourself. But the situation aint great, considering cutting some of the non-essential expenses, and improving your income. Don't worry you are not alone and this is solvable, hand in there."
	SessionDescriptionForecastMillionaire           = "You are now a millionaire. It doesn't matter if your journey is different from the others, you deserve a well-big congratulation here. The first million is always hard to make but you made it at this point, you will be fine from here. However, if this is the point in the far future, maybe you can start improving your income? Reduce your Expense? A dollar more in income or less in expense may go a longer way than you think."
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
