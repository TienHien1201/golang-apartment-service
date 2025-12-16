package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
)

type mailUsecase struct {
	repo repository.MailRepository
}

func NewMailUsecase(repo repository.MailRepository) usecase.MailUsecase {
	return &mailUsecase{repo: repo}
}

func (u *mailUsecase) SendLoginMail(ctx context.Context, email, fullName string) error {
	return u.repo.Send(ctx, repository.MailData{
		Email:    email,
		Subject:  "Cảnh báo đăng nhập",
		Title:    "Cảnh báo đăng nhập",
		Color:    "red",
		Action:   "đăng nhập vào Apartment Business",
		FullName: fullName,
	})
}

func (u *mailUsecase) SendRegisterMail(ctx context.Context, email, fullName string) error {
	return u.repo.Send(ctx, repository.MailData{
		Email:    email,
		Subject:  "Thông báo đăng ký",
		Title:    "Thông báo đăng ký",
		Color:    "green",
		Action:   "đăng ký thành công tài khoản Apartment Business",
		FullName: fullName,
	})
}
