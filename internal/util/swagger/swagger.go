package swagger

import (
	httputil "github.com/M15t/ghoul/pkg/util/http"
	_ "github.com/M15t/ghoul/pkg/util/swagger" // Swagger stuffs
)

// ListRequest holds data of listing request for swagger
// swagger:parameters adminSessionList
type ListRequest struct {
	httputil.ListRequest
}
