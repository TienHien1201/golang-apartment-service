package mail

import (
	"net/smtp"
)

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type Mailer struct {
	auth smtp.Auth
	host string
	addr string
	from string
}

func NewMailer(cfg SMTPConfig, from string) *Mailer {
	auth := smtp.PlainAuth(
		"",
		cfg.User,
		cfg.Pass,
		cfg.Host,
	)

	return &Mailer{
		auth: auth,
		host: cfg.Host,
		addr: cfg.Host + ":" + cfg.Port,
		from: from,
	}
}
