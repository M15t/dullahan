package model

import "github.com/labstack/echo/v4"

// AuthToken holds authentication token details with refresh token
// swagger:model
type AuthToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthAdmin represents data stored in JWT token for admin
type AuthAdmin struct {
	ID   int64
	Name string
	Role string
}

// AuthSession represents data stored in JWT token for session
type AuthSession struct {
	ID   int64
	Code string
	Role string
}

// Auth represents auth interface
type Auth interface {
	Admin(echo.Context) *AuthAdmin
}
