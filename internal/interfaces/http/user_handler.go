package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type ChangeRequest struct {
	NumberID int
}

type ChangeResult struct {
	Status string
}

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (u UserHandler) UpdateDefaultNumber(c echo.Context) error {
	var request ChangeRequest
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invaild user id"})
	}
	err = c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invaild change request"})
	}
	_, err = u.userUsecase.UpdateDefaultNumber(userId, request.NumberID)
	if err != nil {
		switch err.(type) {
		case usecase.InvalidNumber:
			return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
		case usecase.UserNotFound:
			return c.JSON(http.StatusNotFound, Error{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, ChangeResult{Status: "successful"})
}
