package income

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// custome error
var (
	ErrIncomeNotFound = server.NewHTTPError(http.StatusBadRequest, "INCOME_NOTFOUND", "Income not found")
)
