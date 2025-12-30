package permission

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xcontext "thomas.vn/apartment_service/pkg/http/context"
	xlogger "thomas.vn/apartment_service/pkg/logger"

	"thomas.vn/apartment_service/internal/domain/usecase"
)

type PermissionsMiddleware struct {
	logger       *xlogger.Logger
	permissionUC usecase.PermissionUsecase
}

func NewPermissionMiddleware(logger *xlogger.Logger, permissionUC usecase.PermissionUsecase) *PermissionsMiddleware {
	return &PermissionsMiddleware{logger: logger, permissionUC: permissionUC}
}
func (m *PermissionsMiddleware) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get(xcontext.UserContextKey).(*xuser.User)
		if !ok || user == nil {
			return xhttp.UnauthorizedResponse(c, "user not authenticated")
		}

		req := model.CheckPermissionRequest{
			RoleID:   user.RoleID,
			Method:   c.Request().Method,
			Endpoint: c.Path(),
		}

		hasPermission, err := m.permissionUC.CheckPermission(
			c.Request().Context(),
			req,
		)
		if err != nil {
			m.logger.Error(
				"check permission failed",
				xlogger.Error(err),
				xlogger.Int("role_id", user.RoleID),
				xlogger.String("method", req.Method),
				xlogger.String("endpoint", req.Endpoint),
			)

			return xhttp.NewAppError(
				"ERR_INTERNAL_SERVER",
				"",
				"internal server error",
				http.StatusInternalServerError,
			)
		}

		if !hasPermission {
			return xhttp.ForbiddenResponse(c, "permission denied")
		}

		return next(c)
	}
}
