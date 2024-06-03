package debt

// func (s *Debt) recalculateDebt(debt *model.Debt) error {
// 	interestPaidOffEachMonth := s.calculateInterestPaid(debt.AnnualInterest, debt.MonthlyPayment)

// 	return s.db.Debt.Update(s.db.GDB, map[string]interface{}{
// 		"interest_paid_each_month": interestPaidOffEachMonth,
// 		"debt_paid_off_each_month": calculateDebtPaid(debt.MonthlyPayment, interestPaidOffEachMonth),
// 	}, debt.ID)
// }

// func (s *Debt) calculateTotalMonthlyPaymentDebt(sessionID int64) float64 {
// 	var totalMonthlyPaymentDebt float64 = 0.0
// 	s.db.Debt.SumTotalMonthlyPaymentDebt(s.db.GDB, &totalMonthlyPaymentDebt, sessionID)

// 	return totalMonthlyPaymentDebt
// }

// func calculateDebtPaid(monthlyPayment, interestPaid float64) float64 {
// 	return monthlyPayment - interestPaid
// }
