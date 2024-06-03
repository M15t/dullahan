package debt

import (
	"dullahan/internal/db"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New creates new debt application service
func New(db *db.Service, rbacSvc rbac.Intf, cr Crypter) *Debt {
	return &Debt{db: db, rbac: rbacSvc, cr: cr}
}

// Debt represents latefee application service
type Debt struct {
	db   *db.Service
	rbac rbac.Intf
	cr   Crypter
}

// Crypter represents security interface
type Crypter interface {
	RoundFloat(f float64) float64
}
