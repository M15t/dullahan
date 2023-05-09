package income

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"

	structutil "github.com/M15t/ghoul/pkg/util/struct"
)

// Create creates a new latefee account
func (s *Income) Create(c echo.Context, authUsr *model.AuthCustomer, data CreationData) (*model.Income, error) {
	if err := s.enforce(authUsr, model.ActionCreate); err != nil {
		return nil, err
	}

	rec := &model.Income{
		Name:      data.Name,
		Type:      data.Type,
		Amount:    data.Amount,
		SessionID: authUsr.SessionID,
	}

	if err := s.db.Income.Create(s.db.GDB, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error creating latefee").SetInternal(err)
	}

	// * recalcuate total income
	if err := s.updateCurrentSession(authUsr.SessionID); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return rec, nil
}

// Update updates purchase information
func (s *Income) Update(c echo.Context, authUsr *model.AuthCustomer, id int64, data UpdateData) (*model.Income, error) {
	if err := s.enforce(authUsr, model.ActionUpdate); err != nil {
		return nil, err
	}

	// * check legit session
	if existed, err := s.db.Income.Exist(s.db.GDB, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil || !existed {
		return nil, ErrIncomeNotFound
	}

	// optimistic update
	updates := structutil.ToMap(data)
	if err := s.db.Income.Update(s.db.GDB, updates, id); err != nil {
		return nil, server.NewHTTPInternalError("Error updating purchase").SetInternal(err)
	}

	// * get latest record
	rec := new(model.Income)
	if err := s.db.Income.View(s.db.GDB, rec, id); err != nil {
		return nil, ErrIncomeNotFound.SetInternal(err)
	}

	// * recalculate total income
	if err := s.updateCurrentSession(authUsr.SessionID); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return rec, nil
}

// Delete deletes a purchase
func (s *Income) Delete(c echo.Context, authUsr *model.AuthCustomer, id int64) error {
	if err := s.enforce(authUsr, model.ActionDelete); err != nil {
		return err
	}

	if existed, err := s.db.Income.Exist(s.db.GDB, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil || !existed {
		return ErrIncomeNotFound.SetInternal(err)
	}

	if err := s.db.Income.Delete(s.db.GDB, id); err != nil {
		return server.NewHTTPInternalError("Error deleting purchase").SetInternal(err)
	}

	return nil
}

// enforce checks Income permission to perform the action
func (s *Income) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectIncome, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
