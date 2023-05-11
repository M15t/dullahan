package expense

import (
	"dullahan/internal/db"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New creates new expense application service
func New(db *db.Service, rbacSvc rbac.Intf, cr Crypter) *Expense {
	return &Expense{db: db, rbac: rbacSvc, cr: cr}
}

// Expense represents latefee application service
type Expense struct {
	db   *db.Service
	rbac rbac.Intf
	cr   Crypter
}

// Crypter represents security interface
type Crypter interface {
	RoundFloat(f float64) float64
}
