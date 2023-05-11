package session

import (
	"dullahan/internal/model"
)

func (s *Session) calculateTotalIncome(session *model.Session) float64 {
	var totalIncome float64 = 0.0
	if len(session.Incomes) > 0 {
		s.db.Income.SumTotalIncome(s.db.GDB, &totalIncome, session.ID)
	}

	return totalIncome
}

func (s *Session) calculateTotalMonthlyPaymentDebt(session *model.Session) float64 {
	var totalMonthlyPaymentDebt float64 = 0.0
	if len(session.Debts) > 0 {
		s.db.Debt.SumTotalMonthlyPaymentDebt(s.db.GDB, &totalMonthlyPaymentDebt, session.ID)
	}

	return totalMonthlyPaymentDebt
}

func (s *Session) calculateTotalRemaingingDebt(session *model.Session) float64 {
	var totalRemaingingDebt float64 = 0.0
	if len(session.Debts) > 0 {
		s.db.Debt.SumTotalRemainingDebt(s.db.GDB, &totalRemaingingDebt, session.ID)
	}

	return totalRemaingingDebt
}

func (s *Session) recalculateSession(session *model.Session) error {
	var funFund float64 = 0.0
	var isAchivedInvestment bool

	totalEssentialExpense := session.TotalEssentialExpense
	totalNonEssentialExpense := session.TotalNonEssentialExpense
	totalExpense := s.cr.RoundFloat(totalEssentialExpense + totalNonEssentialExpense)

	totalRemaningDebt := s.cr.RoundFloat(s.calculateTotalRemaingingDebt(session))

	totalIncome := s.cr.RoundFloat(s.calculateTotalIncome(session))

	totalMonthlyPaymentDebt := s.cr.RoundFloat(s.calculateTotalMonthlyPaymentDebt(session))

	monthlyNetFlow := totalIncome - (totalMonthlyPaymentDebt + totalExpense)

	emergencyFund := s.cr.RoundFloat(totalEssentialExpense * EmergencyFundRate)
	rainydayFund := s.cr.RoundFloat(totalEssentialExpense * RainydayFundRate)

	retirementPlan := s.cr.RoundFloat(totalEssentialExpense * 12 * RetirementPlanRate)

	isAchivedEmergencyFund := (session.CurrentBalance - totalRemaningDebt) >= emergencyFund
	isAchivedRainydayFund := (session.CurrentBalance - totalRemaningDebt) >= (emergencyFund + rainydayFund)
	isAchivedRetirementPlan := (session.CurrentBalance - totalRemaningDebt) >= retirementPlan

	if isAchivedEmergencyFund && isAchivedRainydayFund {
		funFund = monthlyNetFlow / 100 * 20
		isAchivedInvestment = true
	}

	var status string
	switch {
	case monthlyNetFlow < 0:
		status = model.SessionStatusBD
	case 0 <= monthlyNetFlow && monthlyNetFlow < BankruptCeil:
		status = model.SessionStatusPC2PC
	case BankruptCeil <= monthlyNetFlow && monthlyNetFlow <= totalEssentialExpense:
		status = model.SessionStatusLFF
	case totalEssentialExpense < monthlyNetFlow:
		status = model.SessionStatusGFF
	default:
		status = model.SessionStatusDefault
	}

	// * return latest information
	session.TotalAllIncome = totalIncome
	session.TotalAllExpense = totalExpense
	session.TotalMonthlyPaymentDebt = totalMonthlyPaymentDebt
	session.MonthlyNetFlow = monthlyNetFlow
	session.EmergencyFund = emergencyFund
	session.RainydayFund = rainydayFund
	session.FunFund = funFund
	session.RetirementPlan = retirementPlan
	session.IsAchivedEmergencyFund = isAchivedEmergencyFund
	session.IsAchivedRainydayFund = isAchivedRainydayFund
	session.IsAchivedInvestment = isAchivedInvestment
	session.IsAchivedRetirementPlan = isAchivedRetirementPlan
	session.Status = status
	session.Description = model.SessionStatusDescriptions[status]

	// * update session
	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"total_all_income":           totalIncome,
		"total_all_expense":          totalExpense,
		"total_monthly_payment_debt": totalMonthlyPaymentDebt,
		"monthly_net_flow":           monthlyNetFlow,
		"status":                     status,
		"emergency_fund":             emergencyFund,
		"rainyday_fund":              rainydayFund,
		"fun_fund":                   funFund,
		"retirement_plan":            retirementPlan,
		"is_achived_emergency_fund":  isAchivedEmergencyFund,
		"is_achived_rainyday_fund":   isAchivedRainydayFund,
		"is_achived_investment":      isAchivedInvestment,
		"is_achived_retirement_plan": isAchivedRetirementPlan,
	}, session.ID)
}
