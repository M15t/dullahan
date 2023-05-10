package expense

import (
	"dullahan/internal/model"

	"gorm.io/gorm"
)

func (s *Expense) calculateExpense(sessionID int64, dataType string) float64 {
	// * tricky part
	var newExpense float64 = 0.0
	s.db.Expense.SumExpenseByType(s.db.GDB, &newExpense, dataType, sessionID)

	return newExpense
}

func (s *Expense) updateCurrentSession(sessionID int64, dataType string) error {
	newExpense := s.calculateExpense(sessionID, dataType)
	updates := map[string]interface{}{
		"total_expense": gorm.Expr("total_essential_expense + total_non_essential_expense"),
	}

	switch dataType {
	case model.ExpenseTypeEssential:
		updates["total_essential_expense"] = newExpense
	default:
		updates["total_non_essential_expense"] = newExpense
	}
	// * update session
	return s.db.Session.Update(s.db.GDB, updates, sessionID)
}
