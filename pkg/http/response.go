package xhttp

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DataResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
		Data:    data,
	})
}

func ListResponse(c echo.Context, rows interface{}, total int64) error {
	return DataResponse(c, http.StatusOK, &ListDataResponse{Rows: rows, Total: total})
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
	return DataResponse(c, http.StatusInternalServerError, "Something went wrong")
}

func BadRequestResponse(c echo.Context, data interface{}) error {
	return DataResponse(c, http.StatusBadRequest, data)
}

func UnauthorizedResponse(c echo.Context, data interface{}) error {
	return DataResponse(c, http.StatusUnauthorized, data)
}

func AppErrorResponse(c echo.Context, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return DataResponse(c, appErr.Status, []*AppError{appErr})
	}

	return InternalServerErrorResponse(c)
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
		Data:    &ListDataResponse{Rows: data, Total: total},
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
