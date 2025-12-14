package consts

import xqueue "thomas.vn/apartment_service/pkg/queue"

const (
	// UploadFileJobName is the name of the upload file job
	UploadFileJobName = "upload_file_job"

	// UpdateStaffFileJobName is the name of the update staff file job
	UpdateStaffFileJobName = "update_staff_file_job"

	// UploadFileJobType is the type of the upload file job
	UploadFileJobType = xqueue.MessageType("upload_file_job")

	// ProcessCandidatesJobName is the name of the process candidates job
	ProcessCandidatesJobName = "process_candidates_job"

	// ProcessCandidatesJobType is the type of the process candidates job
	ProcessCandidatesJobType = xqueue.MessageType("process_candidates_job")
)
