package expense

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// Custom error
var (
	ErrExpenseNotFound = server.NewHTTPError(http.StatusBadRequest, "EXPENSE_NOTFOUND", "Expense not found")
)
