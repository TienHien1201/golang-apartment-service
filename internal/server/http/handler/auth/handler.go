package xauth

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger      *xlogger.Logger
	authHandler *AuthHandler
}

type HandlerOption func(*Handler)

func WithAuthUsecase(uc usecase.AuthUsecase) HandlerOption {
	return func(h *Handler) {
		h.authHandler = NewAuthHandler(h.logger, uc)
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
func (h *Handler) Auth() *AuthHandler {
	return h.authHandler
}
