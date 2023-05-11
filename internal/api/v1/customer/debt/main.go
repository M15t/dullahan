package debt

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"

	structutil "github.com/M15t/ghoul/pkg/util/struct"
)

// Create creates a new latefee account
func (s *Debt) Create(c echo.Context, authUsr *model.AuthCustomer, data CreationData) (*model.Debt, error) {
	if err := s.enforce(authUsr, model.ActionCreate); err != nil {
		return nil, err
	}

	interestPaidOffEachMonth := s.calculateInterestPaid(data.AnnualInterest, data.MonthlyPayment)

	rec := &model.Debt{
		Name:                     data.Name,
		RemainingAmount:          data.RemainingAmount,
		MonthlyPayment:           data.MonthlyPayment,
		AnnualInterest:           data.AnnualInterest,
		Type:                     data.Type,
		PaymentDeadline:          datatypes.Date(data.PaymentDeadline),
		InterestPaidOffEachMonth: interestPaidOffEachMonth,
		DebtPaidOffEachMonth:     calculateDebtPaid(data.MonthlyPayment, interestPaidOffEachMonth),
		SessionID:                authUsr.SessionID,
	}

	if err := s.db.Debt.Create(s.db.GDB, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error creating latefee").SetInternal(err)
	}

	return rec, nil
}

// Update updates purchase information
func (s *Debt) Update(c echo.Context, authUsr *model.AuthCustomer, id int64, data UpdateData) (*model.Debt, error) {
	if err := s.enforce(authUsr, model.ActionUpdate); err != nil {
		return nil, err
	}

	// * check legit session
	if existed, err := s.db.Debt.Exist(s.db.GDB, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil || !existed {
		return nil, ErrDebtNotFound
	}

	// optimistic update
	updates := structutil.ToMap(data)
	if err := s.db.Debt.Update(s.db.GDB, updates, id); err != nil {
		return nil, server.NewHTTPInternalError("Error updating purchase").SetInternal(err)
	}

	// * get latest record
	rec := new(model.Debt)
	if err := s.db.Debt.View(s.db.GDB, rec, id); err != nil {
		return nil, ErrDebtNotFound.SetInternal(err)
	}

	// * recalculate debt
	if err := s.recalculateDebt(rec); err != nil {
		return nil, server.NewHTTPInternalError("Error updating current session").SetInternal(err)
	}

	return rec, nil
}

// Delete deletes a purchase
func (s *Debt) Delete(c echo.Context, authUsr *model.AuthCustomer, id int64) error {
	if err := s.enforce(authUsr, model.ActionDelete); err != nil {
		return err
	}

	// * check legit session
	if existed, err := s.db.Debt.Exist(s.db.GDB, `id = ? AND session_id = ?`, id, authUsr.SessionID); err != nil || !existed {
		return ErrDebtNotFound.SetInternal(err)
	}

	if err := s.db.Debt.Delete(s.db.GDB, id); err != nil {
		return server.NewHTTPInternalError("Error deleting purchase").SetInternal(err)
	}

	return nil
}

// enforce checks Debt permission to perform the action
func (s *Debt) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectDebt, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
