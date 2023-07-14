package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type SmsHistoryHandler struct {
	smsHistory *usecase.SmsHistoryUsecase
}

func NewSmsHistoryHandler(smsHistory usecase.SmsHistoryUsecase) SmsHistoryHandler {
	return SmsHistoryHandler{smsHistory: &smsHistory}
}

type SearchSmsHistoryResult struct {
	Items []domain.SMSHistory `json:"items"`
	Count int                 `json:"count"`
}

func (s SmsHistoryHandler) Search(c echo.Context) error {

	words := c.QueryParams()["words"]

	smsHistoryItems, err := s.smsHistory.Search(words)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, SearchSmsHistoryResult{Items: smsHistoryItems, Count: len(smsHistoryItems)})

}
