package xcontext

import (
	"net/http"

	xuser "thomas.vn/apartment_service/internal/domain/model/user"

	"github.com/labstack/echo/v4"
)

const UserContextKey = "user"

func MustGetUser(c echo.Context) (*xuser.User, error) {
	user, ok := c.Get(UserContextKey).(*xuser.User)
	if !ok || user == nil {
		return nil, echo.NewHTTPError(
			http.StatusUnauthorized,
			"User not authenticated",
		)
	}
	return user, nil
}
