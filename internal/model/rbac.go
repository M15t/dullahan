package model

// RBAC roles
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleCustomer   = "customer"
)

// AvailableRoles for validation
// var AvailableRoles = []string{RoleAdmin, RoleCustomer}

// RBAC objects
const (
	ObjectAny     = "*"
	ObjectSession = "session"
	ObjectIncome  = "income"
	ObjectExpense = "expense"
	ObjectDebt    = "debt"
)

// RBAC actions
const (
	ActionAny       = "*"
	ActionViewAll   = "view_all"
	ActionView      = "view"
	ActionCreateAll = "create_all"
	ActionCreate    = "create"
	ActionUpdateAll = "update_all"
	ActionUpdate    = "update"
	ActionDeleteAll = "delete_all"
	ActionDelete    = "delete"
)
