package session

import (
	"net/http"

	"github.com/M15t/ghoul/pkg/server"
)

// custom errors
var (
	ErrSessionNotFound = server.NewHTTPError(http.StatusBadRequest, "SESSION_NOTFOUND", "Session not found")
)
