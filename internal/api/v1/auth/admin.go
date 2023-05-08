package auth

import (
	"dullahan/internal/model"

	"github.com/M15t/ghoul/pkg/server"

	"github.com/labstack/echo/v4"
)

// LoginAdmin logs in the given admin, returns access token
func (s *Auth) LoginAdmin(u *model.AuthAdmin) (*model.AuthToken, error) {
	claims := map[string]interface{}{
		"id":   u.ID,
		"name": u.Name,
		"role": model.RoleAdmin,
	}
	token, expiresin, err := s.jwt.GenerateToken(claims, nil)
	if err != nil {
		return nil, server.NewHTTPInternalError("Error generating token").SetInternal(err)
	}

	refreshToken := s.cr.UID()
	// err = s.adb.Update(s.db, map[string]interface{}{"refresh_token": refreshToken, "last_login": time.Now()}, u.ID)
	// if err != nil {
	// 	return nil, server.NewHTTPInternalError("Error updating admin").SetInternal(err)
	// }

	return &model.AuthToken{AccessToken: token, TokenType: "bearer", ExpiresIn: expiresin, RefreshToken: refreshToken}, nil
}

// AuthenticateAdmin tries to authenticate the admin provided by username and password
func (s *Auth) AuthenticateAdmin(c echo.Context, email, password string) (*model.AuthToken, error) {
	// admin, err := s.adb.FindByEmail(s.db, email)
	// if err != nil || admin == nil || !admin.IsActive {
	// 	return nil, ErrAdminNotExisted.SetInternal(err)
	// }

	// if !s.cr.CompareHashAndPassword(admin.Password, password) {
	// 	return nil, ErrInvalidCredentials
	// }

	// return s.LoginAdmin(admin)

	return nil, nil
}

// RefreshTokenAdmin returns the new access token with expired time extended
func (s *Auth) RefreshTokenAdmin(c echo.Context, data RefreshTokenData) (*model.AuthToken, error) {
	// admin, err := s.adb.FindByRefreshToken(s.db, data.RefreshToken)
	// if err != nil || admin == nil {
	// 	return nil, ErrInvalidRefreshToken.SetInternal(err)
	// }

	// return s.LoginAdmin(admin)

	return nil, nil
}

// Admin returns admin data stored in jwt token
func (s *Auth) Admin(c echo.Context) *model.AuthAdmin {
	id, _ := c.Get("id").(float64)
	name, _ := c.Get("name").(string)
	role, _ := c.Get("role").(string)

	return &model.AuthAdmin{
		ID:   int64(id),
		Name: name,
		Role: role,
	}
}
