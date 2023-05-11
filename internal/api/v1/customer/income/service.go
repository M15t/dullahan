package income

import (
	"dullahan/internal/db"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New creates new income application service
func New(db *db.Service, rbacSvc rbac.Intf, cr Crypter) *Income {
	return &Income{db: db, rbac: rbacSvc, cr: cr}
}

// Income represents latefee application service
type Income struct {
	db   *db.Service
	rbac rbac.Intf
	cr   Crypter
}

// Crypter represents security interface
type Crypter interface {
	RoundFloat(f float64) float64
}
