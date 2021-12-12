package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crypto/rand"
	"errors"
	"math/big"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type DBHandler struct {
	DB *sql.DB
}

func (h *DBHandler) Rotate(ec echo.Context) error {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "jiachengwu"
	var tokenExpiresAt time.Time
	tokenExpiresAt = time.Now().AddDate(0, 0, 1)
	claims["exp"] = tokenExpiresAt.Unix()
	accessToken, tokenErr := token.SignedString([]byte(getLatestSecret(h.DB)))

	if tokenErr != nil {
		return tokenErr
	}

	type response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}

	res := response{
		Message:     "This JWT has expiration of 1 day but active rotation (every 5 minutes) of secret invalidates it in 10 minutes. Rotation keeps two secrets at max",
		AccessToken: accessToken,
	}
	return ec.JSON(http.StatusOK, res)
}

func (h *DBHandler) RotatenAnswer(ec echo.Context) error {

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

	// try check token with db records
	allSecrets := getLatestSecrets(h.DB)
	var latestMessage string
	for _, secret := range allSecrets {
		serr, msg := trySecret(accessToken, secret)
		latestMessage = msg
		if serr == nil {
			res.Message = msg
			return ec.JSON(http.StatusOK, res)
		}
	}

	res.Message = latestMessage
	return ec.JSON(http.StatusUnprocessableEntity, res)

}

func getLatestSecret(db *sql.DB) string {
	var secret string

	row, err := db.Query("SELECT secret FROM secrets ORDER BY id DESC LIMIT 1")
	if err != nil {
		fmt.Println("error selecting from db")
		return ""
	}
	defer row.Close()
	if row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&secret)
		fmt.Println("secret fetched from db")
	} else {
		// no secret so create one
		return insertRandomSecret(db)
	}

	return secret
}

func getLatestSecrets(db *sql.DB) []string {
	var secrets []string

	row, err := db.Query("SELECT secret FROM secrets ORDER BY id DESC LIMIT 2")
	if err != nil {
		fmt.Println("error selecting from db")
		return secrets
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var secret string
		row.Scan(&secret)
		fmt.Println("secret fetched from db")
		secrets = append(secrets, secret)
	}

	return secrets
}

func trySecret(accessToken string, secret string) (error, string) {

	var message string
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if token == nil {
		message = "JWT is invalid or corrupted"
		return errors.New(message), message
	}

	if token.Valid {
		// Access token is completely fine, go ahead to check refresh token
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			message = "JWT has expired"
			return errors.New(message), message
		} else {
			message = "JWT is invalid or corrupted"
			return errors.New(message), message
		}
	} else {
		message = "JWT is invalid or corrupted"
		return errors.New(message), message
	}
	message = "PASSED"
	return nil, message

}

func insertRandomSecret(db *sql.DB) string {
	randomSecretInt, _ := rand.Int(rand.Reader, new(big.Int).SetInt64(1000000000))
	randomSecret := fmt.Sprintf("%x", randomSecretInt)

	insertStudentSQL := `INSERT INTO secrets(secret) VALUES (?)`
	statement, perr := db.Prepare(insertStudentSQL)
	if perr != nil {
		fmt.Errorf("Fatal error, unable to insert to sqlite: %w \n", perr)
	}
	_, ierr := statement.Exec(randomSecret)
	if ierr != nil {
		fmt.Errorf("Fatal error, unable to insert to sqlite: %w \n", ierr)
	}
	return randomSecret
}
