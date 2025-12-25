package xauth

type LoginRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	Token    *string `json:"token"`
}
type AuthInfoResult struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	FullName string `json:"full_name"`
	IsTotp   bool   `json:"isTotp"`
}

type AuthLoginResult struct {
	IsTotp       bool   `json:"isTotp,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
