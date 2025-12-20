package xauth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
)

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	url := h.googleOAuth.AuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Đây là luồng frontend đã xử lí sau khi đang nhập google thành công thì vào thẳng trang chủ của frontend luôn
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
