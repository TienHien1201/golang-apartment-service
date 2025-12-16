package repository

import "context"

type MailData struct {
	Email    string
	Subject  string
	Title    string
	Color    string
	Action   string
	FullName string
}

type MailRepository interface {
	Send(ctx context.Context, data MailData) error
}
