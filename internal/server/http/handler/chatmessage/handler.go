package chatmessage

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger             *xlogger.Logger
	chatMessageHandler *ChatMessagesHandler
}

// # Funtional Options Pattern

type HandlerOption func(*Handler)

// withchatMessageUC ...

func WithChatMessageUsecase(uc usecase.ChatMessageUsecase) HandlerOption {
	return func(h *Handler) {
		h.chatMessageHandler = NewChatMessageHandler(h.logger, uc)
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

// ChatMessage returns the chatmessage handler
func (h *Handler) ChatMessage() *ChatMessagesHandler {
	return h.chatMessageHandler
}
