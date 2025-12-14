package service

import (
	"context"
)

type AIService interface {
	VerifyCV(ctx context.Context, candidateID string, attachFile string, jobCode string, jobDesc string) (string, error)
}
