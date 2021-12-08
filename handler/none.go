package handler

import (
	"net/http"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ExpiredLogin returns an expired JWT that cannot be used, but the secret is easy to crack with brute-force
func NoneLogin(ec echo.Context) error {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "jiachengwu"
	var tokenExpiresAt time.Time
	tokenExpiresAt = time.Now().AddDate(0, 0, -1)
	claims["exp"] = tokenExpiresAt.Unix()
	randomSecret := os.Getenv("randomSecret")
	accessToken, tokenErr := token.SignedString([]byte(randomSecret))

	if tokenErr != nil {
		return tokenErr
	}

	type response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}

	res := response{
		Message:     "This JWT has expired 1 day ago, but you bypass with 'none' algorithm",
		AccessToken: accessToken,
	}
	return ec.JSON(http.StatusOK, res)
}
