package income

func (s *Income) calculateTotalIncome(sessionID int64) float64 {
	// * tricky part
	var totalIncome float64 = 0.0
	s.db.Income.SumTotalIncome(s.db.GDB, &totalIncome, sessionID)

	return totalIncome
}

func (s *Income) updateCurrentSession(sessionID int64) error {
	// * update session
	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"total_income": s.calculateTotalIncome(sessionID),
	}, sessionID)
}
