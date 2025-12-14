package model

type VerifyRequest struct {
	CandidateID int `gorm:"column:candidate_id" param:"candidate_id" json:"candidate_id" swaggerignore:"true" validate:"required,gt=0" example:"11"`
	JobID       int `gorm:"column:job_id" param:"job_id" json:"job_id" swaggerignore:"true" validate:"required,gt=0" example:"11"`
}

type VerifyResponse struct {
	JobTitle            string                `json:"job_title"`
	CandidateEvaluation []CandidateEvaluation `json:"candidate_evaluations"`
	RecruiterSummary    string                `json:"recruiter_summary"`
}

type CandidateEvaluation struct {
	FullName           string   `json:"full_name"`
	Score              int      `json:"score"`
	Strengths          string   `json:"strengths"`
	Weaknesses         string   `json:"weaknesses"`
	OverallFitSummary  string   `json:"overall_fit_summary"`
	InterviewQuestions []string `json:"interview_questions"`
}

type ScanCVRequest struct {
	JobDescription string `json:"jd_text"`
	CVFile         string `json:"files"`
}

type ScanCVResponse struct {
	VerifyResult   int    `json:"verify_result" param:"verify_result"`
	VerifyResponse string `json:"verify_response" param:"verify_response"`
}
