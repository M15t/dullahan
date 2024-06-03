package income

import (
	"dullahan/internal/model"
	"net/http"

	httputil "github.com/M15t/ghoul/pkg/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents income http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents income application interface
type Service interface {
	Create(c echo.Context, authUsr *model.AuthCustomer, data CreationData) (*model.Income, error)
	Update(c echo.Context, authUsr *model.AuthCustomer, id int64, data UpdateData) (*model.Income, error)
	Delete(c echo.Context, authUsr *model.AuthCustomer, id int64) error
}

// NewHTTP creates new card http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/customer/incomes customer-incomes customerIncomeCreate
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

	// swagger:operation PATCH /v1/customer/incomes/{id} customer-incomes customerIncomeUpdate
	// ---
	// summary: Update income information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of income
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerIncomeUpdateData"
	// responses:
	//   "200":
	//     description: The updated income
	//     schema:
	//       "$ref": "#/definitions/Income"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /v1/customer/incomes/{id} customer-incomes customerIncomeDelete
	// ---
	// summary: Deletes an income
	// parameters:
	// - name: id
	//   in: path
	//   description: id of income
	//   type: integer
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.DELETE("/:id", h.delete)
}

// CreationData contains income data from json request
// swagger:model CustomerIncomeCreationData
type CreationData struct {
	// example: Full-Time job earning
	Name string `json:"name" validate:"required,max=100"`
	// example: MONTHLY
	Type string `json:"type" validate:"required,oneof=MONTHLY PASSIVE"`
	// example: 2000
	Amount float64 `json:"amount" validate:"gte=0"`
}

// UpdateData contains income data from json request
// swagger:model CustomerIncomeUpdateData
type UpdateData struct {
	// example: Full-Time job earning
	Name *string `json:"name,omitempty" validate:"omitempty,max=100"`
	// example: MONTHLY
	Type *string `json:"type,omitempty" validate:"omitempty,oneof=MONTHLY PASSIVE"`
	// example: 2000
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,gte=0"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	resp, err := h.svc.Create(c, h.auth.Customer(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqIDint64(c)
	if err != nil {
		return err
	}
	u := UpdateData{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	resp, err := h.svc.Update(c, h.auth.Customer(c), id, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqIDint64(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(c, h.auth.Customer(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
