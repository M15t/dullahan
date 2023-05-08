package rbac

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	// Add permission for admin role
	r.AddPolicy(model.RoleAdmin, model.ObjectAny, model.ActionAny)

	// Roles inheritance
	// r.AddGroupingPolicy(model.RoleSuperAdmin, model.RoleAdmin)

	r.GetModel().PrintPolicy()

	return r
}
