package mail

import (
	"context"
	"fmt"
	"mime"
	"net/smtp"

	"thomas.vn/apartment_service/internal/domain/repository"
)

func (m *Mailer) Send(
	_ context.Context,
	data repository.MailData,
) error {

	subject := mime.QEncoding.Encode("utf-8", data.Subject)
	fromName := mime.QEncoding.Encode("utf-8", "Hien CNTT")

	from := fmt.Sprintf("%s <%s>", fromName, m.fromEmail())

	body := fmt.Sprintf(`
<div>
	<p style="color:%s">%s</p>
	<p>
		Chào <b>%s</b>, %s
	</p>
</div>
`,
		data.Color,
		data.Title,
		data.FullName,
		data.Action,
	)

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		from,
		data.Email,
		subject,
		body,
	)

	return smtp.SendMail(
		m.addr,
		m.auth,
		m.fromEmail(),
		[]string{data.Email},
		[]byte(msg),
	)
}
func (m *Mailer) fromEmail() string {
	return "tienhien.cntt@gmail.com"
}
