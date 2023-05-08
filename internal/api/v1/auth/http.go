package auth

import (
	"net/http"

	"dullahan/internal/model"

	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	svc Service
}

// Service represents auth service interface
type Service interface {
	Start(echo.Context) (*model.AuthToken, error)
	RefreshToken(c echo.Context, data RefreshTokenData) (*model.AuthToken, error)
}

// NewHTTP creates new auth http service
func NewHTTP(svc Service, eg *echo.Group) {
	h := HTTP{svc}

	// swagger:operation POST /v1/start auth authStart
	// ---
	// summary: Start a new session
	// responses:
	//   "200":
	//     description: Access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("/start", h.start)

	// swagger:operation POST /v1/refresh-token auth authRefreshToken
	// ---
	// summary: Refresh access token
	// parameters:
	// - name: token
	//   in: body
	//   description: The given `refresh_token` when login
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RefreshTokenData"
	// responses:
	//   "200":
	//     description: New access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("/refresh-token", h.refreshToken)
}

// RefreshTokenData represents refresh token request data
// swagger:model
type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *HTTP) start(c echo.Context) error {
	resp, err := h.svc.Start(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) refreshToken(c echo.Context) error {
	r := RefreshTokenData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	resp, err := h.svc.RefreshToken(c, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
