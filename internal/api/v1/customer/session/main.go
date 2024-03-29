package session

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/labstack/echo/v4"
)

// Me returns the current session information
func (s *Session) Me(c echo.Context, authUsr *model.AuthCustomer) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionView); err != nil {
		return nil, err
	}

	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB.Preload("Incomes").Preload("Expenses"), rec, authUsr.SessionID); err != nil {
		return nil, ErrSessionNotFound.SetInternal(err)
	}

	return rec, nil
}

// enforce checks Session permission to perform the action
func (s *Session) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectSession, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
