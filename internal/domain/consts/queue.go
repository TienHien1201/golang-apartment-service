package consts

import xqueue "thomas.vn/apartment_service/pkg/queue"

const (
	// Job
	UploadFileJobName       = "upload_file_job"
	MailJobName             = "mail_job"
	UploadUserAvatarJobName = "upload_local_avatar_file_job"

	// Mail payload types
	MailJobType       = xqueue.MessageType("mail_job")
	QueueMailLogin    = xqueue.MessageType("mail_login")
	QueueMailRegister = xqueue.MessageType("mail_register")

	//Upload file local
	UploadUserAvatarJobType = xqueue.MessageType("upload_local_avatar_file_job")
)
