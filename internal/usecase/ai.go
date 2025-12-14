package usecase

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/service"
	"thomas.vn/apartment_service/internal/domain/usecase"
	"thomas.vn/apartment_service/internal/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type aiUsecase struct {
	logger       *xlogger.Logger
	aiRepo       repository.AiRepository
	downloadURL  string
	queueService service.QueueService
}

func NewAiUsecase(logger *xlogger.Logger, aiRepo repository.AiRepository, downloadURL string, queueService service.QueueService) usecase.AiUsecase {
	return &aiUsecase{
		logger:       logger,
		aiRepo:       aiRepo,
		downloadURL:  downloadURL,
		queueService: queueService,
	}
}

func (u *aiUsecase) VerifyCV(attachFile string, jobDesc string) (int, string, string, error) {
	var filePaths []string

	err := json.Unmarshal([]byte(attachFile), &filePaths)
	if err != nil {
		u.logger.Error("Failed to unmarshal attach file", xlogger.Error(err))
		return 0, "", "", err
	}

	filePath := ""

	if len(filePaths) > 0 {
		filePath = strings.ReplaceAll(filePaths[0], "\\/", "/")
		filePath = u.downloadURL + filePath
	}

	verifyResult, VerifyResponse, err := u.aiRepo.VerifyCV(filePath, jobDesc)
	if err != nil {
		u.logger.Error("Failed to verify cv", xlogger.Error(err))
		return 0, "", "", err
	}

	var verifyScore string
	if len(VerifyResponse.CandidateEvaluation) > 0 {
		score := VerifyResponse.CandidateEvaluation[0].Score
		if score >= 65 {
			verifyScore = "2"
		} else {
			verifyScore = "1"
		}
	} else {
		verifyScore = "1"
	}

	jsonResponse, err := json.Marshal(VerifyResponse)
	if err != nil {
		u.logger.Error("Failed to marshal verify response", xlogger.Error(err))
		return 0, "", "", err
	}

	return verifyResult, string(jsonResponse), verifyScore, nil
}

func (u *aiUsecase) UploadCV(attachFile string) string {
	path := "attachments/resumes/" + time.Now().Format("2006-01-02")

	err := u.queueService.PublishMessage(context.Background(), consts.UploadFileJobName, map[string]interface{}{
		"file": attachFile,
		"path": path,
	})

	if err != nil {
		u.logger.Error("Failed to publish message", xlogger.Error(err))
		return ""
	}

	return path
}
