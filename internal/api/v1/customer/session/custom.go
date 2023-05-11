package session

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// custom errors
var (
	ErrSessionNotFound = server.NewHTTPError(http.StatusBadRequest, "SESSION_NOTFOUND", "Session not found")
)

// Const
const (
	BankruptCeil      float64 = 200.00
	EmergencyFundRate         = 6.00
	RainydayFundRate          = 3.00

	RetirementPlanRate = 10.00 // years
)
