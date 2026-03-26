package repository

import (
	"thomas.vn/apartment_service/internal/domain/model"
)

// AiRepository defines the contract for AI/CV verification data access.
// Implementations live in internal/repository and communicate with
// external AI services (HTTP calls, file downloads, etc.).
type AiRepository interface {
	// VerifyCV submits a file URL and job description to the AI service.
	VerifyCV(attachFile string, jobDesc string) (int, model.VerifyResponse, error)

	// VerifyCVDownload downloads the CV file first, then submits to AI service.
	VerifyCVDownload(attachFile string, jobDesc string) (int, model.VerifyResponse, error)
}
