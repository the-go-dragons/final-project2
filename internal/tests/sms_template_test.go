package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/the-go-dragons/final-project2/internal/interfaces/http"
)

func TestNewTemplate(t *testing.T) {
	// Clear the database
	Setup()

	// New request with bad body
	smsTemplate := &handlers.NewSmsTemplateRequest{
		Text: "salam bar to",
	}
	jsonBytes, err := json.Marshal(smsTemplate)
	if err != nil {
		panic(err)
	}

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/new-templates/new", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec, req)
	// Must response with 401 must authenticate
	if rec.Code != http.StatusUnauthorized {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusUnauthorized, rec.Code, rec.Body))
	}

	// Create the user
	user := &handlers.LoginRequest{
		Username: "fazelsamar",
		Password: "123456",
	}
	jsonBytes, err = json.Marshal(user)
	if err != nil {
		panic(err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBytes))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec2.Code, rec2.Body))
	}

	// Login user
	jsonBytes2, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	req3 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBytes2))
	req3.Header.Set("Content-Type", "application/json")
	rec3 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec3.Code, rec3.Body))
	}
	tokenResponse := new(handlers.LoginResponse)
	err = json.Unmarshal(rec3.Body.Bytes(), tokenResponse)
	if err != nil {
		panic(err)
	}

	// request with authentication and bad body
	jsonBytes, err = json.Marshal(smsTemplate)
	if err != nil {
		panic(err)
	}
	req4 := httptest.NewRequest(http.MethodPost, "/new-templates/new", bytes.NewBuffer(jsonBytes))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Authorization", "Bearer "+tokenResponse.Token)
	rec4 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec4, req4)
	// Must response with 403 Must have at least one argument with %s
	if rec4.Code != http.StatusBadRequest {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusBadRequest, rec4.Code, rec4.Body))
	}

	// request with authentication and good body
	smsTemplate = &handlers.NewSmsTemplateRequest{
		Text: "salam bar to %s",
	}
	jsonBytes, err = json.Marshal(smsTemplate)
	if err != nil {
		panic(err)
	}
	req5 := httptest.NewRequest(http.MethodPost, "/new-templates/new", bytes.NewBuffer(jsonBytes))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+tokenResponse.Token)
	rec5 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec5, req5)
	// Must response with 201
	if rec5.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec5.Code, rec5.Body))
	}
}
