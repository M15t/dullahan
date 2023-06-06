package session

import (
	"context"
	"dullahan/internal/model"
	"time"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/allegro/bigcache/v3"
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
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(1*time.Minute))
	defer cache.Close()

	if err := s.calculateSession(cache, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	rec.DataSets = generateDatasets(cache, rec)

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
