package auth

import (
	"dullahan/internal/model"
	"time"

	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"
)

// LoginSession logs in the given session, returns access token
func (s *Auth) LoginSession(session *model.Session) (*model.AuthToken, error) {
	claims := map[string]interface{}{
		"sid":  session.ID,
		"code": session.Code,
		"role": model.RoleCustomer,
	}

	token, expiresin, err := s.jwt.GenerateToken(claims, nil)
	if err != nil {
		return nil, server.NewHTTPInternalError("Error generating token").SetInternal(err)
	}

	refreshToken := s.cr.UID()
	if err := s.db.Session.Update(s.db.GDB, map[string]interface{}{"refresh_token": refreshToken, "last_login": time.Now()}, session.ID); err != nil {
		return nil, server.NewHTTPInternalError("Error updating session").SetInternal(err)
	}

	return &model.AuthToken{AccessToken: token, TokenType: "bearer", ExpiresIn: expiresin, RefreshToken: refreshToken}, nil

}

// Start starts new session
func (s *Auth) Start(c echo.Context) (*model.AuthToken, error) {
	code, err := s.cr.NanoID()
	if err != nil {
		return nil, server.NewHTTPInternalError("Error generating token").SetInternal(err)
	}

	newSession := &model.Session{
		Code:      code,
		IPAddress: c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	if err := s.db.Session.Create(s.db.GDB, newSession); err != nil {
		return nil, server.NewHTTPInternalError("Error creating session").SetInternal(err)
	}

	return s.LoginSession(newSession)
}

// Resume resumes a session
func (s *Auth) Resume(c echo.Context, data CredentialData) (*model.AuthToken, error) {
	rec := new(model.Session)
	if err := s.db.Session.View(s.db.GDB, rec, `code = ?`, data.SessionCode); err != nil {
		return nil, server.NewHTTPInternalError("Error getting session").SetInternal(err)
	}

	return s.LoginSession(rec)
}

// RefreshToken returns the new access token with expired time extended
func (s *Auth) RefreshToken(c echo.Context, data RefreshTokenData) (*model.AuthToken, error) {
	session, err := s.db.Session.FindByRefreshToken(s.db.GDB, data.RefreshToken)
	if err != nil || session == nil {
		return nil, server.NewHTTPInternalError("Error refresh session").SetInternal(err)
	}

	return s.LoginSession(session)
}

// Customer returns customer data stored in jwt token
func (s *Auth) Customer(c echo.Context) *model.AuthCustomer {
	sid, _ := c.Get("sid").(float64)
	code, _ := c.Get("code").(string)
	role, _ := c.Get("role").(string)

	return &model.AuthCustomer{
		SessionID: int64(sid),
		Code:      code,
		Role:      role,
	}
}
