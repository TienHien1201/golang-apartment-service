package handler

import (
	"github.com/labstack/echo/v4"

	"thomas.vn/hr_recruitment/internal/domain/model"
	"thomas.vn/hr_recruitment/internal/domain/usecase"
	xhttp "thomas.vn/hr_recruitment/pkg/http"
	xlogger "thomas.vn/hr_recruitment/pkg/logger"
)

type AiHandler struct {
	logger    *xlogger.Logger
	aiUsecase usecase.AiUsecase
}

func NewAiHandler(logger *xlogger.Logger, aiUsecase usecase.AiUsecase) *AiHandler {
	return &AiHandler{
		logger:    logger,
		aiUsecase: aiUsecase,
	}
}

func (h *AiHandler) VerifyCV(c echo.Context) error {
	var req model.ScanCVRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	// call ai service
	verifyResult, verifyResponse, _, err := h.aiUsecase.VerifyCV(req.CVFile, req.JobDescription)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	response := model.ScanCVResponse{
		VerifyResult:   verifyResult,
		VerifyResponse: verifyResponse,
	}
	return xhttp.SuccessResponse(c, response)
}
