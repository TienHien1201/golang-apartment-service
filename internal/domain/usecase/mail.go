package usecase

import "context"

type MailUsecase interface {
	SendLoginMail(ctx context.Context, email, fullName string) error
	SendRegisterMail(ctx context.Context, email, fullName string) error
}
