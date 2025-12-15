package xmiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/usecase"
)

type PermissionMiddleware struct {
	permissionUC usecase.PermissionUsecase
}

func NewPermissionMiddleware(
	permissionUC usecase.PermissionUsecase,
) *PermissionMiddleware {
	return &PermissionMiddleware{permissionUC: permissionUC}
}

func (m *PermissionMiddleware) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, ok := c.Get(string(UserContextKey)).(*model.User)
		if !ok || user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
		}

		method := c.Request().Method
		endpoint := c.Path()

		hasPermission, err := m.permissionUC.CheckPermission(
			c.Request().Context(),
			user.RoleID,
			method,
			endpoint,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Check permission failed")
		}

		if !hasPermission {
			return echo.NewHTTPError(http.StatusForbidden, "User not permission")
		}

		return next(c)
	}
}
