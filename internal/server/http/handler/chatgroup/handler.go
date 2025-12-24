package chatgroup

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger           *xlogger.Logger
	chatGroupHandler *ChatGroupsHandler
}

type HandlerOption func(*Handler)

func WithChatGroupUsecase(chatGroupUc usecase.ChatGroupUsecase) HandlerOption {
	return func(h *Handler) {
		h.chatGroupHandler = NewChatGroupHandler(h.logger, chatGroupUc)
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

func (h *Handler) ChatGroup() *ChatGroupsHandler {
	return h.chatGroupHandler
}
