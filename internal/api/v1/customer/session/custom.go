package session

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// custom errors
var (
	ErrSessionNotFound = server.NewHTTPError(http.StatusBadRequest, "SESSION_NOTFOUND", "Session not found")

	// Months = []int{12, 24, 36, 48, 60, 72, 84, 96, 108, 120}

	// Quarters = []int{1, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36, 39, 42, 45, 48, 51, 54, 57, 60, 63, 66, 69, 72, 75, 78, 81, 84, 87, 90, 93, 96, 99, 102, 105, 108, 111, 114, 117, 120}
)

// Const
const (
	CustomYear  = 2025
	CustomMonth = 12
	CustomDay   = 31

	UserBudgetIncreasementRate = 0.00966
	BankruptCeil               = 200.00
	EmergencyFundRate          = 6.00
	RainydayFundRate           = 3.00

	RetirementPlanRate = 10.00 // years
)
