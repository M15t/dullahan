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
	updates := map[string]interface{}{}

	switch dataType {
	case model.ExpenseTypeEssential:
		updates["total_essential_expense"] = newExpense
	default:
		updates["total_non_essential_expense"] = newExpense
	}

	// * update session
	// GORM gonna parse the map by alphabetical order then we need to update twice :)
	if err := s.db.Session.Update(s.db.GDB, updates, sessionID); err != nil {
		return err
	}

	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"total_all_expense": gorm.Expr("total_essential_expense + total_non_essential_expense"),
	}, sessionID)
}
