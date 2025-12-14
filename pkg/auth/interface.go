package auth

type TokenVerifier interface {
	VerifyAccessToken(token string) (*Claims, error)
}
