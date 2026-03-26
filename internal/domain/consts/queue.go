package consts

// MessageType is the domain's type alias for queue message types.
// Using a type alias (= string) ensures full compatibility with the
// infrastructure queue implementation without any import coupling.
type MessageType = string

const (
	// Job names (used when registering jobs with the queue)
	UploadFileJobName        = "upload_file_job"
	MailJobName              = "mail_job"
	UploadUserAvatarJobName  = "upload_local_avatar_file_job"
	UploadAvatarCloudJobName = "upload_cloud_avatar_file_job"

	// Message types (used when publishing messages to the queue)
	MailJobType       MessageType = "mail_job"
	QueueMailLogin    MessageType = "mail_login"
	QueueMailRegister MessageType = "mail_register"

	// Upload file local
	UploadUserAvatarJobType MessageType = "upload_local_avatar_file_job"

	// Upload file cloud
	UploadAvatarCloudJobType     MessageType = "upload_cloud_avatar_file_job"
	DeleteCloudinaryAssetJobType MessageType = "delete_cloud_asset_job"
)
