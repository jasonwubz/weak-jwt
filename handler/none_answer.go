package handler

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// NoneAnswer accepts none algorithm used in JWT signing, HINT: header is eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0
func NoneAnswer(ec echo.Context) error {

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
		return jwt.UnsafeAllowNoneSignatureType, nil
	})

	if token == nil {
		res.Message = "JWT is invalid or corrupted"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	if token.Valid {
		// Access token is completely fine, go ahead to check refresh token
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			res.Message = "JWT has expired"
			return ec.JSON(http.StatusUnprocessableEntity, res)
		} else {
			res.Message = "JWT is invalid or corrupted"
			return ec.JSON(http.StatusUnprocessableEntity, res)
		}
	} else {
		res.Message = "JWT is invalid or corrupted"
		return ec.JSON(http.StatusUnprocessableEntity, res)
	}

	res.Message = "Passed"
	return ec.JSON(http.StatusOK, res)

}
