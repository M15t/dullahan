package debt

import (
	"dullahan/internal/model"
	"net/http"
	"time"

	httputil "github.com/M15t/ghoul/pkg/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents debt http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents debt application interface
type Service interface {
	Create(c echo.Context, authUsr *model.AuthCustomer, data CreationData) (*model.Debt, error)
	Update(c echo.Context, authUsr *model.AuthCustomer, id int64, data UpdateData) (*model.Debt, error)
	Delete(c echo.Context, authUsr *model.AuthCustomer, id int64) error
}

// NewHTTP creates new card http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/customer/debts customer-debts customerDebtCreate
	// ---
	// summary: Creates new debt
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerDebtCreationData"
	// responses:
	//   "200":
	//     description: The new debt
	//     schema:
	//       "$ref": "#/definitions/Debt"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation PATCH /v1/customer/debts/{id} customer-debts customerDebtUpdate
	// ---
	// summary: Update debt information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of debt
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerDebtUpdateData"
	// responses:
	//   "200":
	//     description: The updated debt
	//     schema:
	//       "$ref": "#/definitions/Debt"
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

	// swagger:operation DELETE /v1/customer/debts/{id} customer-debts customerDebtDelete
	// ---
	// summary: Deletes an debt
	// parameters:
	// - name: id
	//   in: path
	//   description: id of debt
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

// CreationData contains debt data from json request
// swagger:model CustomerDebtCreationData
type CreationData struct {
	// example: Bank of America
	Name string `json:"name" validate:"required,max=50"`
	//	example: 30000
	RemainingAmount float64 `json:"remaining_amount" validate:"required,gte=0"`
	// example: 500
	MonthlyPayment float64 `json:"monthly_payment" validate:"required,gte=0"`
	// example: 1.1
	AnnualInterest float64 `json:"annual_interest" validate:"required,gte=0"`
	// example: FIXED
	Type string `json:"type" validate:"required,oneof=FIXED FIXED_AMORTIZED FLOAT FLOAT_AMORTIZED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
	// example: 2025-12-31T00:00:00Z
	PaymentDeadline time.Time `json:"payment_deadline"`
}

// UpdateData contains debt data from json request
// swagger:model CustomerDebtUpdateData
type UpdateData struct {
	// example: Bank of America
	Named *string `json:"name,omitempty" validate:"omitempty,max=50"`
	//	example: 30000
	RemainingAmount *float64 `json:"remaining_amount,omitempty" validate:"omitempty,gte=0"`
	// example: 500
	MonthlyPayment *float64 `json:"monthly_payment,omitempty" validate:"omitempty,gte=0"`
	// example: 1.1
	AnnualInterest *float64 `json:"annual_interest,omitempty" validate:"omitempty,gte=0"`
	// example: FIXED
	Type *string `json:"type,omitempty" validate:"omitempty,oneof=FIXED FIXED_AMORTIZED FLOAT FLOAT_AMORTIZED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
	// example: 2025-12-31T00:00:00Z
	PaymentDeadline *time.Time `json:"payment_deadline,omitempty"`
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
