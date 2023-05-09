package session

import (
	"dullahan/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents session http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents session application interface
type Service interface {
	Me(c echo.Context, authUsr *model.AuthCustomer) (*model.Session, error)
}

// NewHTTP creates new card http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation GET /v1/customer/me customer-me customerMe
	// ---
	// summary: Returns current session
	// responses:
	//   "200":
	//     description: Current session
	//     schema:
	//       "$ref": "#/definitions/Session"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("", h.me)
}

func (h *HTTP) me(c echo.Context) error {
	resp, err := h.svc.Me(c, h.auth.Customer(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
