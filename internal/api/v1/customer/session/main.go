package session

import (
	"dullahan/internal/model"
	"fmt"
	"sort"
	"time"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Me returns the current session information
func (s *Session) Me(c echo.Context, authUsr *model.AuthCustomer) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionView); err != nil {
		return nil, err
	}

	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB.Preload("Incomes").Preload("Expenses").Preload("Debts", func(db *gorm.DB) *gorm.DB {
		return db.Order("debts.annual_interest DESC,debts.remaining_amount ASC")
	}), rec, authUsr.SessionID); err != nil {
		return nil, ErrSessionNotFound.SetInternal(err)
	}

	// * recalcuate total income, expense, debt and budget cups
	if err := s.calculateSession(s.cache, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	rec.FullStatus = mappingFullStatus(rec.Status)
	rec.NextNYears = YearsForCalculation

	return rec, nil
}

// Update updates session information
func (s *Session) Update(c echo.Context, authUsr *model.AuthCustomer, data UpdateData) error {
	if err := s.enforce(authUsr, model.ActionUpdate); err != nil {
		return err
	}

	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"current_balance": data.CurrentBalance,
	}, authUsr.SessionID)
}

// GenerateLineChartData generates line chart data
func (s *Session) GenerateLineChartData(c echo.Context, authUsr *model.AuthCustomer) (*LineChartDataResponse, error) {
	if err := s.enforce(authUsr, model.ActionView); err != nil {
		return nil, err
	}

	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB.Preload("Incomes").Preload("Expenses").Preload("Debts", func(db *gorm.DB) *gorm.DB {
		return db.Order("debts.annual_interest DESC,debts.remaining_amount ASC")
	}), rec, authUsr.SessionID); err != nil {
		return nil, ErrSessionNotFound.SetInternal(err)
	}

	linecharts := s.generateLinechart(s.cache, rec)

	s.cache.Reset()

	return &LineChartDataResponse{
		LineCharts: linecharts,
		Debts:      rec.Debts,
	}, nil
}

// GenerateTimelineData generates timeline data
func (s *Session) GenerateTimelineData(c echo.Context, authUsr *model.AuthCustomer) ([]*model.Timeline, error) {
	if err := s.enforce(authUsr, model.ActionView); err != nil {
		return nil, err
	}

	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB.Preload("Incomes").Preload("Expenses").Preload("Debts", func(db *gorm.DB) *gorm.DB {
		return db.Order("debts.annual_interest DESC,debts.remaining_amount ASC")
	}), rec, authUsr.SessionID); err != nil {
		return nil, ErrSessionNotFound.SetInternal(err)
	}

	var timelines []*model.Timeline
	var format = "Jan 2006"

	for _, debt := range rec.Debts {
		dt, _ := time.Parse(format, debt.ForecastPaidOffDate)
		timelines = append(timelines, &model.Timeline{
			Event:       fmt.Sprintf("Debt %s Paid Off", debt.Name),
			Date:        debt.ForecastPaidOffDate,
			Datetime:    dt,
			Description: fmt.Sprintf("Your %s has been paid off. %.2f$ now will be deducted from your expenses", debt.Name, debt.MonthlyPayment),
		})
	}

	if rec.ForecastEmergencyBudgetFilledDate != "" {
		dt, _ := time.Parse(format, rec.ForecastEmergencyBudgetFilledDate)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastEmergencyBudgetFilled,
			Date:        rec.ForecastEmergencyBudgetFilledDate,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastEmergencyBudgetFilled,
		})
	}

	if rec.ForecastRainydayBudgetFilledDate != "" {
		dt, _ := time.Parse(format, rec.ForecastRainydayBudgetFilledDate)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastRainydayBudgetFilled,
			Date:        rec.ForecastRainydayBudgetFilledDate,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastRainydayBudgetFilled,
		})
	}

	if rec.ForecastStartInvestingDate != "" {
		dt, _ := time.Parse(format, rec.ForecastStartInvestingDate)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastStartInvesting,
			Date:        rec.ForecastStartInvestingDate,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastStartInvesting,
		})
	}

	if rec.ForecastFinancialFreedomDate != "" {
		dt, _ := time.Parse(format, rec.ForecastFinancialFreedomDate)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastFinancialFreedom,
			Date:        rec.ForecastFinancialFreedomDate,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastFinancialFreedom,
		})
	}

	if rec.ForecastMillionaireDate != "" {
		dt, _ := time.Parse(format, rec.ForecastMillionaireDate)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastMillionaire,
			Date:        rec.ForecastMillionaireDate,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastMillionaire,
		})
	}

	if rec.ForecastBankrupt != "" {
		dt, _ := time.Parse(format, rec.ForecastBankrupt)
		timelines = append(timelines, &model.Timeline{
			Event:       model.SessionTitleForecastBankrupt,
			Date:        rec.ForecastBankrupt,
			Datetime:    dt,
			Description: model.SessionDescriptionForecastBankrupt,
		})
	}

	sort.Slice(timelines[:], func(i, j int) bool {
		return timelines[i].Datetime.Before(timelines[j].Datetime)
	})

	return timelines, nil
}

// enforce checks Session permission to perform the action
func (s *Session) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectSession, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
