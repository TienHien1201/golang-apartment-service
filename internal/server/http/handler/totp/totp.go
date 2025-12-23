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

func (h *TotpHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
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

func (h *TotpHandler) Save(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.SaveTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	if err := h.uc.Save(ctx, user, req.Secret, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}

func (h *TotpHandler) Verify(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.VerifyTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	if err := h.uc.Verify(ctx, user, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}

func (h *TotpHandler) Disable(c echo.Context) error {
	ctx := c.Request().Context()

	var req dtototp.DisableTotpRequest
	if err := c.Bind(&req); err != nil {
		return xhttp.BadRequestResponse(c, err.Error())
	}

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	if err := h.uc.Disable(ctx, user, req.Token); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, true)
}
