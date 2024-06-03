package main

import (
	"dullahan/config"
	"encoding/json"
	"os"
	"time"

	"dullahan/internal/db"
	"dullahan/internal/model"
	dbutil "dullahan/internal/util/db"

	"gorm.io/datatypes"
)

func main() {
	cfg, err := config.Load()
	checkErr(err)

	gdb, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	checkErr(err)
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	sqlDB, err := gdb.DB()
	defer sqlDB.Close()

	dbSvc := db.New(gdb)

	// * sessions
	if err := migrateSession(dbSvc); err != nil {
		checkErr(err)
	}

	// * income
	if err := migrateIncome(dbSvc); err != nil {
		checkErr(err)
	}

	// * debt
	if err := migrateDebt(dbSvc); err != nil {
		checkErr(err)
	}

	// * expense
	if err := migrateExpense(dbSvc); err != nil {
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type (
	customBool bool

	postgreSQLTimestamp struct {
		time.Time
	}

	postgreSQLDate struct {
		datatypes.Date
	}

	tmpSession struct {
		// The primary key of the record
		ID        int64               `json:"id" gorm:"primary_key"`
		CreatedAt postgreSQLTimestamp `json:"created_at"`
		UpdatedAt postgreSQLTimestamp `json:"updated_at"`
		Code      string              `json:"code" gorm:"type:varchar(20);unique_index"`
		IPAddress string              `json:"ip_address" gorm:"type:varchar(45)" `
		UserAgent string              `json:"user_agent" gorm:"type:text"`

		RefreshToken string               `json:"-" gorm:"type:varchar(100);unique_index"`
		LastLogin    *postgreSQLTimestamp `json:"last_login"`

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

		IsAchivedEmergencyFund  customBool `json:"is_achived_emergency_fund"`
		IsAchivedRainydayFund   customBool `json:"is_achived_rainyday_fund"`
		IsAchivedInvestment     customBool `json:"is_achived_investment"`
		IsAchivedRetirementPlan customBool `json:"is_achived_retirement_plan"`

		ForecastEmergencyBudgetFilledDate string `json:"forecast_emergency_budget_filled_date" gorm:"type:varchar(50)"`
		ForecastRainydayBudgetFilledDate  string `json:"forecast_rainyday_budget_filled_date" gorm:"type:varchar(50)"`
		ForecastStartInvestingDate        string `json:"forecast_start_investing_date" gorm:"type:varchar(50)"`
		ForecastFinancialFreedomDate      string `json:"forecast_financial_freedom_date" gorm:"type:varchar(50)"`
		ForecastMillionaireDate           string `json:"forecast_millionaire_date" gorm:"type:varchar(50)"`
		ForecastBankrupt                  string `json:"forecast_bankrupt" gorm:"type:varchar(50)"`

		Status string `json:"status" gorm:"type:varchar(10)"`
	}

	tmpIncome struct {
		ID        int64               `json:"id" gorm:"primary_key"`
		CreatedAt postgreSQLTimestamp `json:"-"`
		UpdatedAt postgreSQLTimestamp `json:"-"`
		SessionID int64               `json:"session_id"`

		Amount float64 `json:"amount"`
		Name   string  `json:"name" gorm:"type:varchar(100)"`
		Type   string  `json:"type" gorm:"type:varchar(10);default:MONTHLY"` // MONTHLY, PASSIVE
	}

	tmpDebt struct {
		ID        int64               `json:"id" gorm:"primary_key"`
		CreatedAt postgreSQLTimestamp `json:"-"`
		UpdatedAt postgreSQLTimestamp `json:"-"`
		SessionID int64               `json:"session_id"`

		Name            string         `json:"name" gorm:"type:varchar(50)"`
		RemainingAmount float64        `json:"remaining_amount"`
		MonthlyPayment  float64        `json:"monthly_payment"`
		AnnualInterest  float64        `json:"annual_interest"`
		Type            string         `json:"type" gorm:"type:varchar(10);default:FIXED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
		PaymentDeadline postgreSQLDate `json:"payment_deadline"`

		ForecastPaidOffDate string `json:"forecast_paid_off_date" gorm:"type:varchar(50)"`
	}

	tmpExpense struct {
		ID        int64               `json:"id" gorm:"primary_key"`
		CreatedAt postgreSQLTimestamp `json:"-"`
		UpdatedAt postgreSQLTimestamp `json:"-"`
		SessionID int64               `json:"session_id"`

		Amount float64 `json:"amount"`
		Name   string  `json:"name" gorm:"type:varchar(100)"`
		Type   string  `json:"type" gorm:"type:varchar(15);default:ESSENTIAL"` // ESSENTIAL, NON_ESSENTIAL
	}
)

// UnmarshalJSON unmarshals the JSON data into a postgreSQLTimestamp
func (p *postgreSQLTimestamp) UnmarshalJSON(data []byte) error {
	timeStr := string(data)
	parsedTime, err := time.Parse(`"2006-01-02 15:04:05"`, timeStr)
	if err != nil {
		return err
	}
	p.Time = parsedTime
	return nil
}

// UnmarshalJSON unmarshals the JSON data into a postgreSQLDate
func (p *postgreSQLDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	timeStr := string(data)
	parsedDate, err := time.Parse(`"2006-01-02"`, timeStr)
	if err != nil {
		return err
	}

	p.Date = datatypes.Date(parsedDate)
	return nil
}

// UnmarshalJSON customizes the unmarshalling process for customBool
func (cb *customBool) UnmarshalJSON(data []byte) error {
	var numericValue int
	if err := json.Unmarshal(data, &numericValue); err == nil {
		*cb = customBool(numericValue != 0)
		return nil
	}

	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err != nil {
		return err
	}

	*cb = customBool(boolValue)
	return nil
}

func migrateSession(dbSvc *db.Service) error {
	jsonData, err := os.ReadFile("sessions.json")
	if err != nil {
		return err
	}

	var sessions []tmpSession

	// Unmarshal the JSON data into the slice of structs
	err = json.Unmarshal(jsonData, &sessions)
	if err != nil {
		return err
	}

	for _, s := range sessions {
		if err := dbSvc.Session.Create(dbSvc.GDB, &model.Session{
			Base: model.Base{
				ID:        s.ID,
				CreatedAt: s.CreatedAt.Time,
				UpdatedAt: s.UpdatedAt.Time,
			},
			Code:                              s.Code,
			IPAddress:                         s.IPAddress,
			UserAgent:                         s.UserAgent,
			RefreshToken:                      s.RefreshToken,
			LastLogin:                         &s.LastLogin.Time,
			TotalAllIncome:                    s.TotalAllExpense,
			TotalAllExpense:                   s.TotalAllExpense,
			TotalMonthlyPaymentDebt:           s.TotalMonthlyPaymentDebt,
			TotalEssentialExpense:             s.TotalEssentialExpense,
			TotalNonEssentialExpense:          s.TotalNonEssentialExpense,
			MonthlyNetFlow:                    s.MonthlyNetFlow,
			CurrentBalance:                    s.CurrentBalance,
			ActualEmergencyFund:               s.ActualEmergencyFund,
			ActualRainydayFund:                s.ActualRainydayFund,
			ActualFunFund:                     s.ActualFunFund,
			ExpectedEmergencyFund:             s.ExpectedEmergencyFund,
			ExpectedRainydayFund:              s.ActualRainydayFund,
			ExpectedFunFund:                   s.ExpectedFunFund,
			Investment:                        s.Investment,
			RetirementPlan:                    s.RetirementPlan,
			IsAchivedEmergencyFund:            bool(s.IsAchivedEmergencyFund),
			IsAchivedRainydayFund:             bool(s.IsAchivedRainydayFund),
			IsAchivedInvestment:               bool(s.IsAchivedInvestment),
			IsAchivedRetirementPlan:           bool(s.IsAchivedRetirementPlan),
			ForecastEmergencyBudgetFilledDate: s.ForecastEmergencyBudgetFilledDate,
			ForecastRainydayBudgetFilledDate:  s.ForecastRainydayBudgetFilledDate,
			ForecastStartInvestingDate:        s.ForecastStartInvestingDate,
			ForecastFinancialFreedomDate:      s.ForecastFinancialFreedomDate,
			ForecastMillionaireDate:           s.ForecastMillionaireDate,
			ForecastBankrupt:                  s.ForecastBankrupt,
			Status:                            s.Status,
		}); err != nil {
			return err
		}
	}

	return nil
}

func migrateIncome(dbSvc *db.Service) error {
	jsonData, err := os.ReadFile("incomes.json")
	if err != nil {
		return err
	}

	var incomes []tmpIncome

	// Unmarshal the JSON data into the slice of structs
	err = json.Unmarshal(jsonData, &incomes)
	if err != nil {
		return err
	}

	for _, s := range incomes {
		if err := dbSvc.Income.Create(dbSvc.GDB, &model.Income{
			ID:        s.ID,
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
			SessionID: s.SessionID,
			Amount:    s.Amount,
			Name:      s.Name,
			Type:      s.Type,
		}); err != nil {
			return err
		}
	}

	return nil
}

func migrateDebt(dbSvc *db.Service) error {
	jsonData, err := os.ReadFile("debts.json")
	if err != nil {
		return err
	}

	var debts []tmpDebt

	// Unmarshal the JSON data into the slice of structs
	err = json.Unmarshal(jsonData, &debts)
	if err != nil {
		return err
	}

	for _, s := range debts {
		if err := dbSvc.Debt.Create(dbSvc.GDB, &model.Debt{
			ID:                  s.ID,
			CreatedAt:           s.CreatedAt.Time,
			UpdatedAt:           s.UpdatedAt.Time,
			SessionID:           s.SessionID,
			Name:                s.Name,
			RemainingAmount:     s.RemainingAmount,
			MonthlyPayment:      s.MonthlyPayment,
			AnnualInterest:      s.AnnualInterest,
			Type:                s.Type,
			PaymentDeadline:     s.PaymentDeadline.Date,
			ForecastPaidOffDate: s.ForecastPaidOffDate,
		}); err != nil {
			return err
		}
	}

	return nil
}

func migrateExpense(dbSvc *db.Service) error {
	jsonData, err := os.ReadFile("expenses.json")
	if err != nil {
		return err
	}

	var expenses []tmpExpense

	// Unmarshal the JSON data into the slice of structs
	err = json.Unmarshal(jsonData, &expenses)
	if err != nil {
		return err
	}

	for _, s := range expenses {
		if err := dbSvc.Expense.Create(dbSvc.GDB, &model.Expense{
			ID:        s.ID,
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
			SessionID: s.SessionID,
			Amount:    s.Amount,
			Name:      s.Name,
			Type:      s.Type,
		}); err != nil {
			return err
		}
	}

	return nil
}
