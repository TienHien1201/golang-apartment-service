package consts

import xqueue "thomas.vn/apartment_service/pkg/queue"

const (
	// Job
	UploadFileJobName = "upload_file_job"
	MailJobName       = "mail_job"
	MailJobType       = xqueue.MessageType("mail_job")

	// Mail payload types
	QueueMailLogin    = xqueue.MessageType("mail_login")
	QueueMailRegister = xqueue.MessageType("mail_register")
)
