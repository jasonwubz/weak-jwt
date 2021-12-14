package handler

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/martinlindhe/base36"

	"net/http"
	"strings"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// insecureSecretGenerator uses math/rand which is not cryptographically secure
func insecureSecretGenerator() string {
	// this code tries to simulate issue outlined in https://github.com/YMFE/yapi/issues/2117

	// all values of base36
	//'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	//'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	//'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	//'U', 'V', 'W', 'X', 'Y', 'Z'
	//
	rand.Seed(time.Now().UnixNano())
	randFloat := fmt.Sprintf("%v", rand.Float64())[2:]
	fmt.Println(randFloat)

	i, _ := strconv.Atoi(randFloat)

	//output is usually 11 bytes
	return base36.Encode(uint64(i))
}

func MathRand(ec echo.Context) error {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "jiachengwu"
	var tokenExpiresAt time.Time
	tokenExpiresAt = time.Now().AddDate(0, 0, -1)
	claims["exp"] = tokenExpiresAt.Unix()
	randomSecret := os.Getenv("INSECURERANDOMSECRET")

	accessToken, tokenErr := token.SignedString([]byte(randomSecret))

	if tokenErr != nil {
		return tokenErr
	}

	type response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}

	res := response{
		Message:     fmt.Sprintf("This JWT has expired 1 day ago, but it is signed with insecure secret in base 36 represenation of length %d generated using math/rand", len(randomSecret)),
		AccessToken: accessToken,
	}
	return ec.JSON(http.StatusOK, res)
}

func MathRandAnswer(ec echo.Context) error {

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
	randomSecret := os.Getenv("INSECURERANDOMSECRET")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(randomSecret), nil
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
