package permission

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger            *xlogger.Logger
	permissionHandler *PermissionsHandler
}

// # Funtional Options Pattern

type HandlerOption func(*Handler)

// withPermissionUC ...

func WithPermissionUsecase(uc usecase.PermissionUsecase) HandlerOption {
	return func(h *Handler) {
		h.permissionHandler = NewPermissionHandler(h.logger, uc)
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

// Permission returns the Permission handler
func (h *Handler) Permission() *PermissionsHandler {
	return h.permissionHandler
}
