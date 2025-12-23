package xtotp

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger      *xlogger.Logger
	totpHandler *TotpHandler
}

type HandlerOption func(*Handler)

func WithTotpUsecase(uc usecase.TotpUsecase) HandlerOption {
	return func(h *Handler) {
		h.totpHandler = NewTotpHandler(h.logger, uc)
	}
}

func NewHandler(logger *xlogger.Logger, opts ...HandlerOption) *Handler {
	h := &Handler{
		logger: logger,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) Totp() *TotpHandler {
	return h.totpHandler
}
