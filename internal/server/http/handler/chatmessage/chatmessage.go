package chatmessage

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model/chatmessage"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type ChatMessagesHandler struct {
	logger        *xlogger.Logger
	chatMessageUc usecase.ChatMessageUsecase
}

func NewChatMessageHandler(logger *xlogger.Logger, chatMessageUc usecase.ChatMessageUsecase) *ChatMessagesHandler {
	return &ChatMessagesHandler{
		logger:        logger,
		chatMessageUc: chatMessageUc,
	}
}
func (h *ChatMessagesHandler) List(c echo.Context) error {
	var req chatmessage.ListChatMessageRequest

	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	rawFilters := c.QueryParam("filters")
	if rawFilters != "" {
		var f struct {
			ChatGroupID int `json:"chatGroupID"`
		}
		if err := json.Unmarshal([]byte(rawFilters), &f); err != nil {
			return xhttp.BadRequestResponse(c, "filters must be valid JSON")
		}
		req.ChatGroupID = f.ChatGroupID
	}

	if req.ChatGroupID == 0 {
		return xhttp.BadRequestResponse(c, "ChatGroupID is required")
	}

	res, total, err := h.chatMessageUc.ListChatMessages(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("List chat messages failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.PaginationListResponse(c, &req.PaginationOptions, res, total)
}
