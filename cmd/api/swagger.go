//go:build swagger

// Dullahan - Calculating personal finance
//
// # API documents for project
//
// `Believe that you are worthy of financial freedom. Do something you love and then all you ever have to do is be yourself to succeed - Jen Sincero.`
//
// ## Authentication
// All API endpoints started with version, ex: `/v1/customer/*`, require authentication token.
// Firstly, grab the **access_token** from the response of `/start`. Then include this header in all API calls:
// ```
// Authorization: Bearer ${access_token}
// ```
//
// Terms Of Service: N/A
//
// Version: 1.0.0
// Contact: Nguyen Nguyen <khanhnguyen1411@gmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearer: []
//
// SecurityDefinitions:
// bearer:
//   type: apiKey
//   name: Authorization
//   in: header
//
// swagger:meta
package main

import (
	"embed"
)

//go:embed swagger-ui
var embedSwaggerui embed.FS

func init() {
	enableSwagger = true
	swaggerui = embedSwaggerui
}

// OK - No Content
// swagger:response ok
type swaggOKResp struct{}