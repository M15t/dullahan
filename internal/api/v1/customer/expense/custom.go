package expense

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// custome error
var (
	ErrExpenseNotFound = server.NewHTTPError(http.StatusBadRequest, "EXPENSE_NOTFOUND", "Expense not found")
)
