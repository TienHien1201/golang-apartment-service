package xauth

import (
	"strings"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"

	"thomas.vn/apartment_service/internal/domain/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	logger   *xlogger.Logger
	tokenUc  usecase.TokenUsecase
	userRepo repository.UserRepository
}

func NewAuthMiddleware(
	logger *xlogger.Logger,
	tokenUsecase usecase.TokenUsecase,
	userRepo repository.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger,
		tokenUc:  tokenUsecase,
		userRepo: userRepo,
	}
}
func (m *AuthMiddleware) Protect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return xhttp.UnauthorizedResponse(c, "authorization header is missing")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return xhttp.UnauthorizedResponse(c, "invalid authorization header format")
		}

		accessToken := strings.TrimSpace(parts[1])
		if accessToken == "" {
			return xhttp.UnauthorizedResponse(c, "access token is empty")
		}

		claims, err := m.tokenUc.VerifyAccessToken(accessToken)
		if err != nil {
			m.logger.Warn(
				"verify access token failed",
				xlogger.Error(err),
			)
			return xhttp.UnauthorizedResponse(c, "invalid or expired access token")
		}

		user, err := m.userRepo.GetUserByID(c.Request().Context(), claims.UserID)
		if err != nil {
			m.logger.Error(
				"get user by id failed",
				xlogger.Error(err),
			)
			return xhttp.InternalServerErrorResponse(c)
		}

		if user == nil {
			return xhttp.UnauthorizedResponse(c, "user not found or inactive")
		}

		c.Set(string(UserContextKey), user)

		return next(c)
	}
}
