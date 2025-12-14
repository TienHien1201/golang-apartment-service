package xuser

import (
	"github.com/labstack/echo/v4"

	"thomas.vn/hr_recruitment/internal/domain/model"
	"thomas.vn/hr_recruitment/internal/domain/usecase"
	xhttp "thomas.vn/hr_recruitment/pkg/http"
	xlogger "thomas.vn/hr_recruitment/pkg/logger"
)

type UserHandler struct {
	logger *xlogger.Logger
	userUC usecase.UserUsecase
}

func NewUserHandler(logger *xlogger.Logger, userUC usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		logger: logger,
		userUC: userUC,
	}
}

// Create godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Param data body model.CreateUserRequest true "Create user request"
// @Success 201 {object} xhttp.APIResponse{data=model.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/v2/users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var req model.CreateUserRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, err := h.userUC.CreateUser(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Create user failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.CreatedResponse(c, res)
}

// Get godoc
// @Summary Get user
// @Description Get user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} xhttp.APIResponse{data=model.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/v2/users/{id} [get]
func (h *UserHandler) Get(c echo.Context) error {
	var req model.UserIDRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, err := h.userUC.GetUser(c.Request().Context(), req.ID)
	if err != nil {
		h.logger.Error("Get user failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, res)
}

// Update godoc
// @Summary Update user
// @Description Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body model.UpdateUserRequest true "Update user request"
// @Success 200 {object} xhttp.APIResponse{data=model.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/v2/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	var req model.UpdateUserRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, err := h.userUC.UpdateUser(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Update user failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, res)
}

// Delete godoc
// @Summary Delete user
// @Description Delete user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} xhttp.APIResponse{}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/v2/users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	var req model.UserIDRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	err := h.userUC.DeleteUser(c.Request().Context(), req.ID)
	if err != nil {
		h.logger.Error("Delete user failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, nil)
}

// List godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param from_date query string false "From date"
// @Param to_date query string false "To date"
// @Param range_by query string false "Range by"
// @Param order_by query string false "Order by"
// @Param sort_by query string false "Sort by"
// @Param status query int false "Status"
// @Success 200 {object} xhttp.APIResponse{data=model.ListUserResponse}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/v2/users [get]
func (h *UserHandler) List(c echo.Context) error {
	var req model.ListUserRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	res, total, err := h.userUC.ListUsers(c.Request().Context(), &req)
	if err != nil {
		h.logger.Error("List users failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.ListResponse(c, res, total)
}
