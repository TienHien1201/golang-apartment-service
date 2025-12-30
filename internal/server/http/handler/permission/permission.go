package permission

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xcontext "thomas.vn/apartment_service/pkg/http/context"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type PermissionsHandler struct {
	logger       *xlogger.Logger
	permissionUC usecase.PermissionUsecase
}

func NewPermissionHandler(logger *xlogger.Logger, permissionUC usecase.PermissionUsecase) *PermissionsHandler {
	return &PermissionsHandler{logger: logger, permissionUC: permissionUC}
}

func (h *PermissionsHandler) Create(c echo.Context) error {
	var req model.CreatePermissionRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	user, ok := c.Get(xcontext.UserContextKey).(*xuser.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	result, err := h.permissionUC.CreatePermission(
		c.Request().Context(),
		&req,
		user.ID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}
