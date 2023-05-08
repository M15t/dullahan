package income

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
)

// Create creates a new latefee account
func (s *Income) Create(authUsr *model.AuthCustomer, data CreationData) (*model.Income, error) {
	if err := s.enforce(authUsr, model.ActionCreate); err != nil {
		return nil, err
	}

	rec := &model.Income{}

	if err := s.db.Income.Create(s.db.GDB, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error creating latefee").SetInternal(err)
	}

	return rec, nil
}

// enforce checks Income permission to perform the action
func (s *Income) enforce(authUsr *model.AuthCustomer, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectIncome, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
