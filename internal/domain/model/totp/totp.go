package dtototp

type SaveTotpRequest struct {
	Secret string `json:"secret" validate:"required"`
	Token  string `json:"token" validate:"required"`
}

type VerifyTotpRequest struct {
	Token string `json:"token" validate:"required"`
}

type DisableTotpRequest struct {
	Token string `json:"token" validate:"required"`
}
