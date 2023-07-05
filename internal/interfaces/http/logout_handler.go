package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
)

func (uh *UserHandler) Logout(c echo.Context) error {
	user := c.Get("user").(domain.User)

	// update IsLoginRequired field in user
	user.IsLoginRequired = true
	uh.userUsecase.Update(&user)

	return c.JSON(http.StatusOK, Response{Message: "you logged out successfully"})
}
