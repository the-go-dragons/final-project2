package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"net/http"
	"strconv"
)

type PhoneBookHandler struct {
	phonebookService *usecase.PhoneBookService
}

func NewPhoneBookHandler(phonebookService usecase.PhoneBookService) PhoneBookHandler {
	return PhoneBookHandler{phonebookService: &phonebookService}
}

func (n PhoneBookHandler) Create(c echo.Context) error {
	var req usecase.PhoneBookDto
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	if len(req.Name) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid name"})
	}

	if req.UserID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid userId"})
	}

	dto := usecase.PhoneBookDto{
		Name:        req.Name,
		UserID:      req.UserID,
		Description: req.Description,
	}

	_, err = n.phonebookService.Create(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n PhoneBookHandler) Edit(c echo.Context) error {
	var req usecase.PhoneBookDto
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	if len(req.Name) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid name"})
	}

	if req.UserID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid userId"})
	}

	if req.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid Id"})
	}

	dto := usecase.PhoneBookDto{
		Name:        req.Name,
		UserID:      req.UserID,
		Description: req.Description,
		ID:          req.ID,
	}

	_, err = n.phonebookService.Edit(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n PhoneBookHandler) GetByUserName(c echo.Context) error {
	username := c.QueryParam("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid username"})
	}

	phonebookList, err := n.phonebookService.GetByUserName(username)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, phonebookList)
}

func (n PhoneBookHandler) GetAll(c echo.Context) error {
	phonebookList, err := n.phonebookService.GetAll()

	if phonebookList != nil && len(phonebookList) == 0 {
		return c.JSON(http.StatusOK, phonebookList)
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, phonebookList)
}

func (n PhoneBookHandler) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "0" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid id"})
	}

	iId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}
	err = n.phonebookService.Delete(uint(iId))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "PhoneBook Deleted Successfully"})
}
