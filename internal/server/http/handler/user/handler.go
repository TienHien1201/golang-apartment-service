package xuser

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger      *xlogger.Logger
	userHandler *UserHandler
}

// # Funtional Options Pattern

type HandlerOption func(*Handler)

// withuserUC ...

func WithUserUsecase(uc usecase.UserUsecase) HandlerOption {
	return func(h *Handler) {
		h.userHandler = NewUserHandler(h.logger, uc)
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

// User returns the user handler
func (h *Handler) User() *UserHandler {
	return h.userHandler
}
