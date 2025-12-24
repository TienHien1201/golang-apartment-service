package chatgroup

import (
	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model/chatgroup"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type ChatGroupsHandler struct {
	logger      *xlogger.Logger
	ChatGroupUC usecase.ChatGroupUsecase
}

func NewChatGroupHandler(logger *xlogger.Logger, chatGroupUc usecase.ChatGroupUsecase) *ChatGroupsHandler {
	return &ChatGroupsHandler{
		logger:      logger,
		ChatGroupUC: chatGroupUc,
	}
}

func (h *ChatGroupsHandler) List(c echo.Context) error {
	var req chatgroup.ListChatGroupRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, total, err := h.ChatGroupUC.ListChatGroups(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("List chat messages failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.PaginationListResponse(c, &req.PaginationOptions, res, total)
}
