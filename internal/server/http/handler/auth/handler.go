package xauth

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xgoogle "thomas.vn/apartment_service/pkg/oauth/google"
)

type Handler struct {
	logger      *xlogger.Logger
	authUC      usecase.AuthUsecase
	googleOAuth *xgoogle.Client
	authHandler *AuthHandler
}

type HandlerOption func(*Handler)

func WithAuthUsecase(uc usecase.AuthUsecase) HandlerOption {
	return func(h *Handler) {
		h.authUC = uc
	}
}

func WithGoogleOAuth(client *xgoogle.Client) HandlerOption {
	return func(h *Handler) {
		h.googleOAuth = client
	}
}

func NewHandler(logger *xlogger.Logger, opts ...HandlerOption) *Handler {
	h := &Handler{
		logger: logger,
	}

	for _, opt := range opts {
		opt(h)
	}

	h.authHandler = NewAuthHandler(
		h.logger,
		h.authUC,
		h.googleOAuth,
	)

	return h
}

func (h *Handler) Auth() *AuthHandler {
	return h.authHandler
}
