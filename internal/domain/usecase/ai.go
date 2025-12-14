package usecase

type AiUsecase interface {
	VerifyCV(attachFile string, jobDesc string) (int, string, string, error)
	UploadCV(attachFile string) string
}
