package xmiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/pkg/auth"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	logger   *xlogger.Logger
	tokenSvc auth.TokenVerifier
	userRepo repository.UserRepository
}

func NewAuthMiddleware(
	logger *xlogger.Logger,
	tokenSvc auth.TokenVerifier,
	userRepo repository.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger,
		tokenSvc: tokenSvc,
		userRepo: userRepo,
	}
}

func (m *AuthMiddleware) Protect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Authorization")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token type")
		}

		accessToken := parts[1]
		if accessToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Access Token")
		}

		claims, err := m.tokenSvc.VerifyAccessToken(accessToken)
		if err != nil {
			m.logger.Warn("Verify access token failed", xlogger.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
		}

		user, err := m.userRepo.GetUserByID(c.Request().Context(), claims.UserID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Get user failed")
		}
		if user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not User")
		}

		c.Set(string(UserContextKey), user)

		return next(c)
	}
}
