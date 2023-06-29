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

func TestUserRegister(t *testing.T) {
	// Clear the database
	Setup()
	// New user and parse it to json
	user := &handlers.SignupRequest{
		Username: "fazelsamar",
		Password: "123456",
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Perform the request
	RouteApp.E.ServeHTTP(rec, req)

	// Check the response status code
	if rec.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec.Code, rec.Body))
	}

	// Repeat the same user and expect 409 user already exists with the given username
	jsonBytes2, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBytes2))
	req2.Header.Set("Content-Type", "application/json")

	rec2 := httptest.NewRecorder()

	RouteApp.E.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusConflict {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusConflict, rec2.Code, rec2.Body))
	}
}

func TestUserLoginAndLogout(t *testing.T) {
	// Create the user
	user := &handlers.LoginRequest{
		Username: "fazelsamar2",
		Password: "123456",
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec, req)

	// Login user
	jsonBytes2, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBytes2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec2.Code, rec2.Body))
	}

	// Logout user and check token
	tokenResponse := new(handlers.LoginResponse)
	err = json.Unmarshal(rec2.Body.Bytes(), tokenResponse)
	if err != nil {
		panic(err)
	}
	req3 := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req3.Header.Set("Authorization", "Bearer "+tokenResponse.Token)
	rec3 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusOK, rec3.Code, rec3.Body))
	}

	// Logout again with same token
	req4 := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req4.Header.Set("Authorization", "Bearer "+tokenResponse.Token)
	rec4 := httptest.NewRecorder()
	RouteApp.E.ServeHTTP(rec4, req4)

	if rec4.Code != http.StatusUnauthorized {
		panic(fmt.Sprintf("Expected status code %d but got %d. The body: %s", http.StatusUnauthorized, rec4.Code, rec4.Body))
	}
}
