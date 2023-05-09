package rbac

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	// Add permisssion for customer role
	r.AddPolicy(model.RoleCustomer, model.ObjectSession, model.ActionView)

	r.AddPolicy(model.RoleCustomer, model.ObjectIncome, model.ActionCreate)
	r.AddPolicy(model.RoleCustomer, model.ObjectIncome, model.ActionUpdate)
	r.AddPolicy(model.RoleCustomer, model.ObjectIncome, model.ActionDelete)

	r.AddPolicy(model.RoleCustomer, model.ObjectExpense, model.ActionCreate)
	r.AddPolicy(model.RoleCustomer, model.ObjectExpense, model.ActionUpdate)
	r.AddPolicy(model.RoleCustomer, model.ObjectExpense, model.ActionDelete)

	r.AddPolicy(model.RoleCustomer, model.ObjectDebt, model.ActionCreate)
	r.AddPolicy(model.RoleCustomer, model.ObjectDebt, model.ActionUpdate)
	r.AddPolicy(model.RoleCustomer, model.ObjectDebt, model.ActionDelete)

	// Add permission for admin role
	r.AddPolicy(model.RoleAdmin, model.ObjectAny, model.ActionAny)

	// Roles inheritance
	// r.AddGroupingPolicy(model.RoleSuperAdmin, model.RoleAdmin)

	r.GetModel().PrintPolicy()

	return r
}
