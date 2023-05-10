package expense

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"

	structutil "github.com/M15t/ghoul/pkg/util/struct"
)

// Create creates a new latefee account
func (s *Expense) Create(c echo.Context, authUsr *model.AuthCustomer, data CreationData) (*model.Expense, error) {
	if err := s.enforce(authUsr, model.ActionCreate); err != nil {
		return nil, err
	}

	rec := &model.Expense{
		Name:      data.Name,
		Type:      data.Type,
		Amount:    data.Amount,
		SessionID: authUsr.SessionID,
	}

	if err := s.db.Expense.Create(s.db.GDB, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error creating latefee").SetInternal(err)
	}

	// * recalcuate total expense
	if err := s.updateCurrentSession(authUsr.SessionID, data.Type); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return rec, nil
}

// Update updates purchase information
func (s *Expense) Update(c echo.Context, authUsr *model.AuthCustomer, id int64, data UpdateData) (*model.Expense, error) {
	if err := s.enforce(authUsr, model.ActionUpdate); err != nil {
		return nil, err
	}

	// * check legit session
	if existed, err := s.db.Expense.Exist(s.db.GDB, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil || !existed {
		return nil, ErrExpenseNotFound
	}

	// optimistic update
	updates := structutil.ToMap(data)
	if err := s.db.Expense.Update(s.db.GDB, updates, id); err != nil {
		return nil, server.NewHTTPInternalError("Error updating purchase").SetInternal(err)
	}

	// * get latest record
	rec := new(model.Expense)
	if err := s.db.Expense.View(s.db.GDB, rec, id); err != nil {
		return nil, ErrExpenseNotFound.SetInternal(err)
	}

	// * recalculate total expense
	if err := s.updateCurrentSession(authUsr.SessionID, rec.Type); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return rec, nil
}

// Delete deletes a purchase
func (s *Expense) Delete(c echo.Context, authUsr *model.AuthCustomer, id int64) error {
	if err := s.enforce(authUsr, model.ActionDelete); err != nil {
		return err
	}

	// * check legit session
	rec := new(model.Expense)
	if err := s.db.Expense.View(s.db.GDB, rec, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil {
		return ErrExpenseNotFound.SetInternal(err)
	}

	if err := s.db.Expense.Delete(s.db.GDB, id); err != nil {
		return server.NewHTTPInternalError("Error deleting purchase").SetInternal(err)
	}

	// * recalculate total expense
	if err := s.updateCurrentSession(authUsr.SessionID, rec.Type); err != nil {
		return server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return nil
}

// enforce checks Expense permission to perform the action
func (s *Expense) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectExpense, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
