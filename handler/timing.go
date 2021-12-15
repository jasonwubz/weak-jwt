package handler

import (
	"fmt"
	"time"

	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func Timing(ec echo.Context) error {

	token := jwt.New(SigningMethodCS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "jiachengwu"
	var tokenExpiresAt time.Time
	tokenExpiresAt = time.Now().AddDate(0, 0, -1)
	claims["exp"] = tokenExpiresAt.Unix()
	randomSecret := "imhardcodedbutyoudontneedtoknowme"

	accessToken, tokenErr := token.SignedString([]byte(randomSecret))

	if tokenErr != nil {
		return tokenErr
	}

	type response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}

	res := response{
		Message:     "This JWT has expired 1 day ago, but its verification uses unsafe string comparison of the hashes",
		AccessToken: accessToken,
	}
	return ec.JSON(http.StatusOK, res)
}

func TimingAnswer(ec echo.Context) error {

	type response struct {
		Message string `json:"message"`
	}

	res := response{}

	authVal := ec.Request().Header.Get("Authorization")

	if len(authVal) == 0 {
		res.Message = "JWT is missing"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	tokenParts := strings.Split(authVal, " ")
	if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "bearer") || len(tokenParts[1]) == 0 {
		res.Message = "JWT is invalid or corrupted"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	accessToken := tokenParts[1]

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("imhardcodedbutyoudontneedtoknowme"), nil
	})

	if token == nil {
		res.Message = "JWT is invalid or corrupted, empty"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	if token.Valid {
		// Access token is completely fine
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			res.Message = "JWT has expired"
			return ec.JSON(http.StatusUnprocessableEntity, res)
		} else {
			res.Message = "JWT is invalid or corrupted"
			//fmt.Println(err)
			return ec.JSON(http.StatusUnprocessableEntity, res)
		}
	} else {
		//fmt.Println(err)
		res.Message = "JWT is invalid or corrupted"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	res.Message = "Passed"
	return ec.JSON(http.StatusOK, res)

}
