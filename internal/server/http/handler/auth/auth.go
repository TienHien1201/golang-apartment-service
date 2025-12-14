package xauth

import (
	"github.com/labstack/echo/v4"

	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type AuthHandler struct {
	logger *xlogger.Logger
	authUC usecase.AuthUsecase
}

func NewAuthHandler(logger *xlogger.Logger, authUC usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		authUC: authUC,
	}
}

// REGISTER
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.CreateUserRequest

	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	user, err := h.authUC.Register(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Register failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.CreatedResponse(c, user)
}

// LOGIN
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req loginRequest

	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	accessToken, refreshToken, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Login failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// REFRESH TOKEN
type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	ctx := c.Request().Context()
	var req refreshRequest

	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	newAccessToken, err := h.authUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error("Refresh token failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"access_token": newAccessToken,
	})
}

// LOGOUT
func (h *AuthHandler) Logout(c echo.Context) error {
	err := h.authUC.Logout(c.Request().Context())

	if err != nil {
		h.logger.Error("Logout failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"message": "Logged out successfully",
	})
}
