package xhttp

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/pkg/query"
)

func DataResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
		Data:    data,
	})
}

func ListResponse(c echo.Context, rows interface{}, total int64) error {
	return DataResponse(c, http.StatusOK, &ListDataResponse{
		Rows:  rows,
		Total: total,
	})
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return DataResponse(c, http.StatusOK, data)
}

func CreatedResponse(c echo.Context, data interface{}) error {
	return DataResponse(c, http.StatusCreated, data)
}

func NoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func InternalServerErrorResponse(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
		Data:    nil,
	})
}

func BadRequestResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusBadRequest, APIResponse{
		Status:  http.StatusBadRequest,
		Message: "Bad Request",
		Data:    data,
	})
}

func ForbiddenResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusForbidden, APIResponse{
		Status:  http.StatusForbidden,
		Message: "Forbidden",
		Data:    data,
	})
}

func UnauthorizedResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusUnauthorized, APIResponse{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Data:    data,
	})
}

func AppErrorResponse(c echo.Context, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return c.JSON(appErr.Status, APIResponse{
			Status:  appErr.Status,
			Message: appErr.Message,
			Data:    nil,
		})
	}

	return InternalServerErrorResponse(c)
}

func PaginationListResponse(
	c echo.Context,
	req *query.PaginationOptions,
	items interface{},
	total int64,
) error {

	page := req.Page
	limit := req.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var totalPage int64 = 1
	if !req.ExcludeTotal && limit > 0 {
		totalPage = (total + int64(limit) - 1) / int64(limit)
	}

	return DataResponse(c, http.StatusOK, &PaginationResponse{
		Page:      page,
		PageSize:  limit,
		TotalItem: total,
		TotalPage: totalPage,
		Items:     items,
	})
}

func OldSuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, OldAPIResponse{
		Message: "Success",
		Status:  1,
		Code:    200,
		Data:    data,
	})
}

func OldListSuccessResponse(c echo.Context, data interface{}, total int64) error {
	return c.JSON(http.StatusOK, OldAPIResponse{
		Message: "Success",
		Status:  1,
		Code:    200,
		Data: &ListDataResponse{
			Rows:  data,
			Total: total,
		},
	})
}

func OldErrorResponse(c echo.Context, message string, code int, data interface{}) error {
	return c.JSON(http.StatusOK, OldAPIResponse{
		Message: message,
		Status:  0,
		Code:    code,
		Data:    data,
	})
}

func OldBadRequestResponse(c echo.Context, data interface{}) error {
	return OldErrorResponse(c, "Bad Request", 400, data)
}

func OldInternalErrorResponse(c echo.Context) error {
	return OldErrorResponse(c, "Internal Server Error", 500, "Something went wrong")
}
