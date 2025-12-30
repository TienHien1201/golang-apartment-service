package xuser

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	user2 "thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xcontext "thomas.vn/apartment_service/pkg/http/context"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type UserHandler struct {
	logger *xlogger.Logger
	userUC user2.UserUsecase
}

func NewUserHandler(logger *xlogger.Logger, userUC user2.UserUsecase) *UserHandler {
	return &UserHandler{
		logger: logger,
		userUC: userUC,
	}
}

// Create godoc
// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
// @Param data body xuser.CreateUserRequest true "Create user request"
// @Success 201 {object} xhttp.APIResponse{data=xuser.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var req xuser.CreateUserRequest
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
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} xhttp.APIResponse{data=xuser.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 404 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users/{id} [get]
func (h *UserHandler) Get(c echo.Context) error {
	var req xuser.UserIDRequest
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
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body xuser.UpdateUserRequest true "Update user request"
// @Success 200 {object} xhttp.APIResponse{data=xuser.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	var req xuser.UpdateUserRequest
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
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} xhttp.APIResponse{}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	var req xuser.UserIDRequest
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
// @Description List users with pagination and filters
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param sort_by query string false "Sort by field"
// @Param order_by query string false "Order by asc/desc"
// @Param filters query string false "JSON encoded filters"
// @Success 200 {object} xhttp.APIResponse{data=[]xuser.User}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users [get]
func (h *UserHandler) List(c echo.Context) error {
	var req xuser.ListUserRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	var filters xuser.UserFilters

	if req.Filters != "" && req.Filters != "undefined" {
		decoded, err := url.QueryUnescape(req.Filters)
		if err != nil {
			return xhttp.BadRequestResponse(c, err)
		}

		if err := json.Unmarshal([]byte(decoded), &filters); err != nil {
			return xhttp.BadRequestResponse(c, err)
		}
	}

	res, total, err := h.userUC.ListUsers(
		c.Request().Context(),
		&req,
		&filters,
	)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}
	return xhttp.PaginationListResponse(c, &req.PaginationOptions, res, total)
}

// UploadLocal godoc
// @Summary Upload avatar (local)
// @Description Upload avatar image to local storage
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} xhttp.APIResponse
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users/avatar-local [post]
func (h *UserHandler) UploadLocal(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return xhttp.NewAppError(
			"ERR_FILE_NOT_FOUND",
			"avatar",
			"Avatar file is required",
			http.StatusBadRequest,
		)
	}
	if err := xhttp.ValidateImageFile(file, 5<<20); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	if err := h.userUC.UploadLocal(ctx, &xuser.UploadAvatarLocalRequest{
		UserID: uint(user.ID),
		File:   file,
	}); err != nil {
		h.logger.Error("Upload local file failed", xlogger.Error(err))
		return xhttp.AppErrorResponse(c, err)
	}
	return xhttp.SuccessResponse(c, map[string]string{
		"message": "Upload avatar local success",
	})
}

// UploadCloud godoc
// @Summary Upload avatar (cloud)
// @Description Upload avatar image to cloud storage
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} xhttp.APIResponse
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 401 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/users/avatar-cloud [post]
func (h *UserHandler) UploadCloud(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := xcontext.MustGetUser(c)
	if err != nil {
		return xhttp.BadRequestResponse(c, err)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return xhttp.NewAppError(
			"ERR_FILE_NOT_FOUND",
			"avatar",
			"Avatar file is required",
			http.StatusBadRequest,
		)
	}

	if err := xhttp.ValidateImageFile(file, 5<<20); err != nil {
		return xhttp.AppErrorResponse(c, err)
	}

	if err := h.userUC.UploadCloud(ctx, &xuser.UploadAvatarCloudRequest{
		UserID: uint(user.ID),
		File:   file,
	}); err != nil {
		h.logger.Error(
			"Upload cloud avatar failed",
			xlogger.Error(err),
			xlogger.Uint("user_id", uint(user.ID)),
		)
		return xhttp.AppErrorResponse(c, err)
	}

	return xhttp.SuccessResponse(c, map[string]string{
		"message": "Avatar upload cloudinary is being processed",
	})
}
