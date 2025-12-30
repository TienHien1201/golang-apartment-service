package xauth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
)

// GoogleLogin godoc
// @Summary Google OAuth login
// @Description Redirect user to Google OAuth consent screen
// @Tags auth
// @Produce json
// @Success 302 "Redirect to Google login page"
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Router /api/auth/google [get]
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	url := h.googleOAuth.AuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handle Google OAuth callback and redirect to frontend with access & refresh tokens
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 302 "Redirect to frontend login callback with tokens"
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	ctx := c.Request().Context()

	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing code")
	}

	token, err := h.googleOAuth.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	profile, err := h.googleOAuth.GetProfile(ctx, token)
	if err != nil {
		return err
	}

	gUser := &model.GoogleUser{
		Email:         profile.Email,
		EmailVerified: profile.EmailVerified,
		FullName:      profile.Name,
		Avatar:        profile.Picture,
		GoogleID:      profile.ID,
	}

	accessToken, refreshToken, err := h.authUC.GoogleLogin(ctx, gUser)
	if err != nil {
		return err
	}

	redirectURL := "http://localhost:3000/login-callback" +
		"?accessToken=" + accessToken +
		"&refreshToken=" + refreshToken

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
