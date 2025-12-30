package xauth

import (
	"github.com/labstack/echo/v4"
	xauth "thomas.vn/apartment_service/internal/domain/model/auth"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xcontext "thomas.vn/apartment_service/pkg/http/context"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xgoogle "thomas.vn/apartment_service/pkg/oauth/google"
)

type AuthHandler struct {
	logger      *xlogger.Logger
	authUC      usecase.AuthUsecase
	googleOAuth *xgoogle.Client
}

func NewAuthHandler(logger *xlogger.Logger, authUC usecase.AuthUsecase, googleOAuth *xgoogle.Client) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authUC:      authUC,
		googleOAuth: googleOAuth,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body xuser.CreateUserRequest true "Register request"
// @Success 201 {object} xhttp.APIResponse{data=xuser.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req xuser.CreateUserRequest

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

// Login godoc
// @Summary Login
// @Description Login with email & password (support TOTP)
// @Tags auth
// @Accept json
// @Produce json
// @Param data body xauth.LoginRequest true "Login request"
// @Success 200 {object} xhttp.APIResponse{data=xauth.AuthLoginResult}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req xauth.LoginRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	result, err := h.authUC.Login(ctx, req.Email, req.Password, req.Token)
	if err != nil {
		h.logger.Error("Login failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	if result.IsTotp {
		return xhttp.SuccessResponse(c, map[string]any{
			"isTotp":       true,
			"accessToken":  result.AccessToken,
			"refreshToken": result.RefreshToken,
		})
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param data body xauth.RefreshRequest true "Refresh token request"
// @Success 200 {object} xhttp.APIResponse{}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	ctx := c.Request().Context()
	var req xauth.RefreshRequest

	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	newAccessToken, newRefreshToken, err := h.authUC.RefreshToken(ctx, req.AccessToken, req.RefreshToken)
	if err != nil {
		h.logger.Error("Refresh token failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

// Logout godoc
// @Summary Logout
// @Description Logout current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} xhttp.APIResponse{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Router /api/auth/logout [post]
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

// GetInfo godoc
// @Summary Get current user info
// @Description Get information of current authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} xhttp.APIResponse{data=xauth.AuthInfoResult}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/auth/get-info [get]
func (h *AuthHandler) GetInfo(c echo.Context) error {
	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return err
	}

	result, err := h.authUC.GetInfo(c.Request().Context(), user)
	if err != nil {
		return err
	}

	return xhttp.SuccessResponse(c, result)
}
