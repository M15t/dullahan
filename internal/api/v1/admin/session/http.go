package session

import (
	"net/http"

	dbutil "github.com/M15t/ghoul/pkg/util/db"

	httputil "github.com/M15t/ghoul/pkg/util/http"

	"dullahan/internal/model"

	"github.com/labstack/echo/v4"
)

// HTTP represents session http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents session application interface
type Service interface {
	Create(authUsr *model.AuthAdmin, data CreationData) (*model.Session, error)
	View(authUsr *model.AuthAdmin, id int) (*model.Session, error)
	List(authUsr *model.AuthAdmin, lq *dbutil.ListQueryCondition, count *int64) ([]*model.Session, error)
	Update(authUsr *model.AuthAdmin, id int, data UpdateData) (*model.Session, error)
	Delete(authUsr *model.AuthAdmin, id int) error
}

// NewHTTP creates new session http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/admin/sessions admin-sessions adminSessionCreate
	// ---
	// summary: Creates new session
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/AdminSessionCreationData"
	// responses:
	//   "200":
	//     description: The new session
	//     schema:
	//       "$ref": "#/definitions/Session"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation GET /v1/admin/sessions/{id} admin-sessions adminSessionView
	// ---
	// summary: Returns a single session
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: The session
	//     schema:
	//       "$ref": "#/definitions/Session"
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
	eg.GET("/:id", h.view)

	// swagger:operation GET /v1/admin/sessions admin-sessions adminSessionList
	// ---
	// summary: Returns list of sessions
	// responses:
	//   "200":
	//     description: List of sessions
	//     schema:
	//       "$ref": "#/definitions/AdminSessionListResp"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("", h.list)

	// swagger:operation PATCH /v1/admin/sessions/{id} admin-sessions adminSessionUpdate
	// ---
	// summary: Update session's information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/AdminSessionUpdateData"
	// responses:
	//   "200":
	//     description: The updated session
	//     schema:
	//       "$ref": "#/definitions/Session"
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

	// swagger:operation DELETE /v1/admin/sessions/{id} admin-sessions adminSessionDelete
	// ---
	// summary: Deletes a session
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
	//   type: integer
	//   required: true
	// responses:
	//   "200":
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

// ListResp contains list of sessions and current page number response
// swagger:model AdminSessionListResp
type ListResp struct {
	Data       []*model.Session `json:"data"`
	TotalCount int64            `json:"total_count"`
}

// CreationData contains session data from json request
// swagger:model AdminSessionCreationData
type CreationData struct {
	IPAddress string `json:"ip_address" validate:"required,max=45"`
	UserAgent string `json:"user_agent" validate:"required"`
}

// UpdateData contains session data from json request
// swagger:model AdminSessionUpdateData
type UpdateData struct {
	IPAddress *string `json:"ip_address" validate:"required,max=45"`
	UserAgent *string `json:"user_agent" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	resp, err := h.svc.Create(h.auth.Admin(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.View(h.auth.Admin(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httputil.ReqListQuery(c)
	if err != nil {
		return err
	}
	var count int64 = 0
	resp, err := h.svc.List(h.auth.Admin(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListResp{resp, count})
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	u := UpdateData{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	usr, err := h.svc.Update(h.auth.Admin(c), id, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(h.auth.Admin(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
