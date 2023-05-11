package debt

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// Custom error
var (
	ErrDebtNotFound = server.NewHTTPError(http.StatusBadRequest, "DEBT_NOTFOUND", "Debt not found")
)

// Const
const (
	BankruptCeil float64 = 200.00
)
