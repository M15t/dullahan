package expense

import (
	"dullahan/internal/model"
)

func (s *Expense) calculateExpense(sessionID int64, dataType string) float64 {
	// * tricky part
	var newExpense float64 = 0.0
	s.db.Expense.SumExpenseByType(s.db.GDB, &newExpense, dataType, sessionID)

	return newExpense
}

func (s *Expense) updateCurrentSession(sessionID int64, dataType string) error {
	// * just recalculate the total of each type in sessions tbl
	newExpense := s.cr.RoundFloat(s.calculateExpense(sessionID, dataType))
	updates := map[string]interface{}{}

	switch dataType {
	case model.ExpenseTypeEssential:
		updates["total_essential_expense"] = newExpense
	default:
		updates["total_non_essential_expense"] = newExpense
	}

	// * update session
	return s.db.Session.Update(s.db.GDB, updates, sessionID)
}
