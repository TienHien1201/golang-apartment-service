package permission

import (
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

// Create godoc
// @Summary Create permission
// @Description Create a new permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param data body model.CreatePermissionRequest true "Create permission request"
// @Success 201 {object} xhttp.APIResponse{data=model.Permission}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/permissions [post]
func (h *PermissionsHandler) Create(c echo.Context) error {
	var req model.CreatePermissionRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	user, ok := c.Get(xcontext.UserContextKey).(*xuser.User)
	if !ok || user == nil {
		return xhttp.BadRequestResponse(c, "Unauthorized")
	}
	result, err := h.permissionUC.CreatePermission(
		c.Request().Context(),
		&req,
		user.ID,
	)
	if err != nil {
		return err
	}

	return xhttp.CreatedResponse(c, result)
}

// Get godoc
// @Summary Get permission
// @Description Get permission by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} xhttp.APIResponse{data=model.Permission}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 404 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/permissions/{id} [get]
func (h *PermissionsHandler) Get(c echo.Context) error {
	var req model.PermissionIDRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, err := h.permissionUC.GetPermissionByID(c.Request().Context(), req.ID)
	if err != nil {
		h.logger.Error("Get permission failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, res)
}

// Update godoc
// @Summary Update permission
// @Description Update permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param data body model.UpdatePermissionRequest true "Update permission request"
// @Success 200 {object} xhttp.APIResponse{data=model.Permission}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 404 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Security BearerAuth
// @Router /api/permissions/{id} [put]
func (h *PermissionsHandler) Update(c echo.Context) error {
	var req model.UpdatePermissionRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}
	res, err := h.permissionUC.UpdatePermission(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Update permission failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}
	return xhttp.SuccessResponse(c, res)

}
