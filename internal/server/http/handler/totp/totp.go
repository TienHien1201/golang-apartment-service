package xtotp

import (
	"github.com/labstack/echo/v4"
	dtototp "thomas.vn/apartment_service/internal/domain/model/totp"
	utotp "thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xcontext "thomas.vn/apartment_service/pkg/http/context"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type TotpHandler struct {
	logger *xlogger.Logger
	uc     utotp.TotpUsecase
}

func NewTotpHandler(logger *xlogger.Logger, uc utotp.TotpUsecase) *TotpHandler {
	return &TotpHandler{
		logger: logger,
		uc:     uc,
	}
}

// Generate godoc
// @Summary Generate TOTP secret
// @Description Generate TOTP secret and QR code for current user
// @Tags totp
// @Accept json
// @Produce json
// @Success 200 {object} xhttp.APIResponse{data=map[string]string}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/totp/generate [post]
func (h *TotpHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	secret, qrCode, err := h.uc.Generate(ctx, user)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"secret": secret,
		"qrCode": qrCode,
	})
}

// Save godoc
// @Summary Save TOTP secret
// @Description Save and enable TOTP after verifying token
// @Tags totp
// @Accept json
// @Produce json
// @Param data body dtototp.SaveTotpRequest true "Save TOTP request"
// @Success 200 {object} xhttp.APIResponse{data=boolean}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/totp/save [post]
func (h *TotpHandler) Save(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.SaveTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	if err := h.uc.Save(ctx, user, req.Secret, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}

// Verify godoc
// @Summary Verify TOTP token
// @Description Verify TOTP token for current user
// @Tags totp
// @Accept json
// @Produce json
// @Param data body dtototp.VerifyTotpRequest true "Verify TOTP request"
// @Success 200 {object} xhttp.APIResponse{data=boolean}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/totp/verify [post]
func (h *TotpHandler) Verify(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.VerifyTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	if err := h.uc.Verify(ctx, user, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}

// Disable godoc
// @Summary Disable TOTP
// @Description Disable TOTP for current user
// @Tags totp
// @Accept json
// @Produce json
// @Param data body dtototp.DisableTotpRequest true "Disable TOTP request"
// @Success 200 {object} xhttp.APIResponse{data=boolean}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/totp/disable [post]
func (h *TotpHandler) Disable(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.DisableTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	if err := h.uc.Disable(ctx, user, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}
