package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
)

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(domain.User)
		if user.ID <= 0 {
			return c.JSON(http.StatusBadRequest, MassageResponse{Message: "User is not found"})
		}
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, MassageResponse{Message: "User is not admin"})
		}
		return next(c)
	}
}
