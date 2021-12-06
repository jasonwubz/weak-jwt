package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// ExpiredLogin returns an expired JWT that cannot be used, but the secret is easy to crack with brute-force
func ExpiredLogin(ec echo.Context) error {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "jiachengwu"
	var tokenExpiresAt time.Time
	tokenExpiresAt = time.Now().AddDate(0, 0, -1)
	claims["exp"] = tokenExpiresAt.Unix()
	accessToken, tokenErr := token.SignedString([]byte(viper.GetString("secret")))

	if tokenErr != nil {
		return tokenErr
	}

	type response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}

	res := response{
		Message:     "This JWT is expired 1 day ago, but you can brute-force the secret and forge your own signature",
		AccessToken: accessToken,
	}
	return ec.JSON(http.StatusOK, res)
}
