package model

type GoogleUser struct {
	RoleID        int `json:"role_id"`
	Email         string
	EmailVerified bool
	FullName      string
	Avatar        string
	GoogleID      string
}
