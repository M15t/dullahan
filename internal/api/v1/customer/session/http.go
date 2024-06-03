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
	Update(c echo.Context, authUsr *model.AuthCustomer, data UpdateData) error
	GenerateLineChartData(c echo.Context, authUsr *model.AuthCustomer) (*LineChartDataResponse, error)
	GenerateTimelineData(c echo.Context, authUsr *model.AuthCustomer) ([]*model.Timeline, error)
}

// NewHTTP creates new card http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation GET /v1/customer/me customer-me customerMe
	// ---
	// summary: Return current session
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

	// swagger:operation GET /v1/customer/me/generate-line-chart customer-me customerMeGenerateLineChart
	// ---
	// summary: Generate line chart data
	// responses:
	//   "200":
	//     description: Line chart data
	//     schema:
	//       "$ref": "#/definitions/LineChartDataResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("/generate-line-chart", h.generateLineChart)

	// swagger:operation GET /v1/customer/me/generate-timeline-chart customer-me customerMeGenerateTimelineChart
	// ---
	// summary: Generate timeline chart data
	// responses:
	//   "200":
	//     description: Timeline chart data
	//     schema:
	//       "$ref": "#/definitions/TimelineChartDataResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("/generate-timeline-chart", h.generateTimelineChart)

	// swagger:operation PATCH /v1/customer/me customer-me customerMeUpdate
	// ---
	// summary: Update current session
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerMeUpdateData"
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
	eg.PATCH("", h.update)
}

// UpdateData contains session data from json request
// swagger:model CustomerMeUpdateData
type UpdateData struct {
	// example: 10000
	CurrentBalance float64 `json:"current_balance" validate:"required,gte=0"`
}

// LineChartDataResponse contains line chart data
// swagger:model
type LineChartDataResponse struct {
	LineCharts []*model.LineChart `json:"data"`
	Debts      []*model.Debt      `json:"debts,omitempty"`
}

// TimelineChartDataResponse contains timeline chart data
// swagger:model
type TimelineChartDataResponse struct {
	Timelines []*model.Timeline `json:"data"`
}

func (h *HTTP) me(c echo.Context) error {
	resp, err := h.svc.Me(c, h.auth.Customer(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) generateLineChart(c echo.Context) error {
	resp, err := h.svc.GenerateLineChartData(c, h.auth.Customer(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LineChartDataResponse{LineCharts: resp.LineCharts, Debts: resp.Debts})
}

func (h *HTTP) generateTimelineChart(c echo.Context) error {
	resp, err := h.svc.GenerateTimelineData(c, h.auth.Customer(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TimelineChartDataResponse{Timelines: resp})
}

func (h *HTTP) update(c echo.Context) error {
	u := UpdateData{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	if err := h.svc.Update(c, h.auth.Customer(c), u); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
