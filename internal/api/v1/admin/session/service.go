package session

import (
	"dullahan/internal/db"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New creates new Session application service
func New(db *db.Service, rbac rbac.Intf) *Session {
	return &Session{db: db, rbac: rbac}
}

// Session represents Session application service
type Session struct {
	db   *db.Service
	rbac rbac.Intf
}
