package session

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"
)

// Me returns the current session information
func (s *Session) Me(c echo.Context, authUsr *model.AuthCustomer) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionView); err != nil {
		return nil, err
	}

	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB.Preload("Incomes").Preload("Expenses").Preload("Debts"), rec, authUsr.SessionID); err != nil {
		return nil, ErrSessionNotFound.SetInternal(err)
	}

	// * recalcuate total income, expense, debt and budget cups
	if err := s.recalculateSession(rec); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	rec.TotalAssets = calculateTotalAssets(rec)

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

// enforce checks Session permission to perform the action
func (s *Session) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectSession, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
