package articles

import (
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Handler struct {
	logger        *xlogger.Logger
	articlHandler *ArticleHandler
}

type HandlerOption func(*Handler)

func WithArticleUsecase(uc usecase.ArticlesUsecase) HandlerOption {
	return func(h *Handler) {
		h.articlHandler = NewArticleHandler(h.logger, uc)
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

func (h *Handler) Articles() *ArticleHandler {
	return h.articlHandler
}
