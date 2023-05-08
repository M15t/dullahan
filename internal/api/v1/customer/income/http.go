package income

import (
	"dullahan/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents income http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents income application interface
type Service interface {
	Create(authUsr *model.AuthCustomer, data CreationData) (*model.Income, error)
}

// NewHTTP creates new card http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/customer/income customer-income customerIncomeCreate
	// ---
	// summary: Creates new income
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerIncomeCreationData"
	// responses:
	//   "200":
	//     description: The new income
	//     schema:
	//       "$ref": "#/definitions/Income"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)
}

// CreationData contains income data from json request
// swagger:model CustomerIncomeCreationData
type CreationData struct {
	// example: 1
	Name string `json:"name" validate:"required"`
	// example: MONTHLY
	Type string `json:"type" validate:"required,oneof=MONTHLY PASSIVE"`
	// example: 1
	Amount float64 `json:"amount" validate:"required,gte=0"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	resp, err := h.svc.Create(h.auth.Customer(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
