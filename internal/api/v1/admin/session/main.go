package session

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
	dbutil "github.com/M15t/ghoul/pkg/util/db"
)

// Create creates a new session
func (s *Session) Create(authUsr *model.AuthAdmin, data CreationData) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionCreateAll); err != nil {
		return nil, err
	}

	session := new(model.Session)

	return session, nil
}

// View returns single session
func (s *Session) View(authUsr *model.AuthAdmin, id int) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionViewAll); err != nil {
		return nil, err
	}

	session := new(model.Session)

	return session, nil
}

// List returns list of sessions
func (s *Session) List(authUsr *model.AuthAdmin, lq *dbutil.ListQueryCondition, count *int64) ([]*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionViewAll); err != nil {
		return nil, err
	}

	sessions := make([]*model.Session, 0)

	return sessions, nil
}

// Update updates session information
func (s *Session) Update(authUsr *model.AuthAdmin, id int, data UpdateData) (*model.Session, error) {
	if err := s.enforce(authUsr, model.ActionUpdateAll); err != nil {
		return nil, err
	}

	session := new(model.Session)

	return session, nil
}

// Delete deletes a session
func (s *Session) Delete(authUsr *model.AuthAdmin, id int) error {
	if err := s.enforce(authUsr, model.ActionDeleteAll); err != nil {
		return err
	}

	return nil
}

// enforce checks permission to perform the action
func (s *Session) enforce(authUsr *model.AuthAdmin, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectSession, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
