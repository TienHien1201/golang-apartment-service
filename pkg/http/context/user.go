package xcontext

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
)

const UserContextKey = "user"

func MustGetUser(c echo.Context) (*model.User, error) {
	user, ok := c.Get(UserContextKey).(*model.User)
	if !ok || user == nil {
		return nil, echo.NewHTTPError(
			http.StatusUnauthorized,
			"User not authenticated",
		)
	}
	return user, nil
}
