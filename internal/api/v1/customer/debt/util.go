package debt

import (
	"dullahan/internal/model"
)

func (s *Debt) recalculateDebt(debt *model.Debt) error {
	interestPaidOffEachMonth := s.calculateInterestPaid(debt.AnnualInterest, debt.MonthlyPayment)

	return s.db.Debt.Update(s.db.GDB, map[string]interface{}{
		"interest_paid_each_month": interestPaidOffEachMonth,
		"debt_paid_off_each_month": calculateDebtPaid(debt.MonthlyPayment, interestPaidOffEachMonth),
	}, debt.ID)
}

func (s *Debt) calculateTotalMonthlyPaymentDebt(sessionID int64) float64 {
	// * tricky part
	var totalMonthlyPaymentDebt float64 = 0.0
	s.db.Debt.SumTotalMonthlyPaymentDebt(s.db.GDB, &totalMonthlyPaymentDebt, sessionID)

	return totalMonthlyPaymentDebt
}

func (s *Debt) calculateTotalIncome(sessionID int64) float64 {
	// * tricky part
	var totalIncome float64 = 0.0
	s.db.Income.SumTotalIncome(s.db.GDB, &totalIncome, sessionID)

	return totalIncome
}

func (s *Debt) calculateExpense(sessionID int64) float64 {
	// * tricky part
	var totalExpense float64 = 0.0
	s.db.Expense.SumExpense(s.db.GDB, &totalExpense, sessionID)

	return totalExpense
}

func (s *Debt) calculateEssentialExpense(sessionID int64) float64 {
	// * tricky part
	var totalEssentialExpense float64 = 0.0
	s.db.Expense.SumExpenseByType(s.db.GDB, &totalEssentialExpense, model.ExpenseTypeEssential, sessionID)

	return totalEssentialExpense
}

func (s *Debt) updateCurrentSession(sessionID int64) error {
	totalMonthlyPaymentDebt := s.calculateTotalMonthlyPaymentDebt(sessionID)
	totalIncome := s.calculateTotalIncome(sessionID)
	totalExpense := s.calculateExpense(sessionID)
	totalEssentialExpense := s.calculateEssentialExpense(sessionID)

	monthlyNetFlow := totalIncome - (totalMonthlyPaymentDebt + totalExpense)

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

	// * update session
	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"total_monthly_payment_debt": totalMonthlyPaymentDebt,
		"monthly_net_flow":           monthlyNetFlow,
		"status":                     status,
	}, sessionID)
}

func (s *Debt) calculateInterestPaid(annualInterest, monthlyPayment float64) float64 {
	return s.cr.RoundFloat((annualInterest / 12.0) * monthlyPayment)
}

func calculateDebtPaid(monthlyPayment, interestPaid float64) float64 {
	return monthlyPayment - interestPaid
}
