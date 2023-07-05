package http

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/config"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func GenerateToken(user *domain.User) (string, error) {
	expirationHoursCofig := config.Config.Jwt.Token.Expire.Hours
	JwtTokenSecretConfig := config.Config.Jwt.Token.Secret.Key

	duration := time.Duration(expirationHoursCofig) * time.Hour
	expirationTime := time.Now().Add(duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    expirationTime.Unix(),
	})

	secretKey := []byte(JwtTokenSecretConfig)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uh *UserHandler) Login(c echo.Context) error {
	var request LoginRequest
	var user *domain.User

	// Check the body data
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid body request"})
		// TODO: all Responses should be in a standard

		// TODO: response messages should be mutli language
		// we can use i18 library
	}

	if request.Username == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}

	// Check for dupplication
	user, err = uh.userUsecase.GetUserByUsername(request.Username)
	if err != nil {
		return c.JSON(http.StatusConflict, Response{Message: "No user found with this credentials"})
	}

	// Check if password is correct
	equalErr := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password))

	if equalErr == nil {
		token, err := GenerateToken(user)

		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{Message: "Server Error"})
		}

		// update IsLoginRequired field
		user.IsLoginRequired = false
		uh.userUsecase.Update(user)
		SetUserToSession(c, user)

		return c.JSON(http.StatusOK, LoginResponse{Message: "You logged in successfully", Token: token})
	}

	return c.JSON(http.StatusConflict, Response{Message: "No user found with this credentials"})
}

func SetUserToSession(c echo.Context, user *domain.User) {
	session := c.Get("session").(*sessions.Session)
	session.Values["userID"] = user.ID
	session.Save(c.Request(), c.Response())
}
