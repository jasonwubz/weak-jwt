package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"weak-jwt/handler"

	"github.com/labstack/echo/v4"
	"github.com/martinlindhe/base36"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config not found: %w \n", err))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Error loading config: %w \n", err))
		}
	}

	//create sqlite file if not exists
	createOpenSqlite()

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	if hasSecretsTable(sqliteDatabase) == false {
		fmt.Println("table doesn't exist")
		createTable(sqliteDatabase)
	}

	dbHandler := handler.DBHandler{
		DB: sqliteDatabase,
	}

	go func() {
		for {
			time.Sleep(time.Duration(300) * time.Second)
			rotateSecrets(sqliteDatabase)
		}
	}()

	ec := echo.New()
	ec.HideBanner = true

	// generate some random passwords
	randomSecret := generateRandomSecret()
	os.Setenv("RANDOMSECRET", randomSecret)
	os.Setenv("INSECURERANDOMSECRET", insecureSecretGenerator())

	ec.POST("/api/weak", handler.ExpiredLogin)
	ec.POST("/api/weak-answer", handler.ExpiredLoginAnswer)
	ec.POST("/api/math-rand", handler.MathRand)
	ec.POST("/api/math-rand-answer", handler.MathRandAnswer)
	ec.POST("/api/none", handler.NoneLogin)
	ec.POST("/api/none-answer", handler.NoneAnswer)
	ec.POST("/api/rotate", dbHandler.Rotate)
	ec.POST("/api/rotate-answer", dbHandler.RotatenAnswer)
	ec.POST("/api/timing", handler.Timing)
	ec.POST("/api/timing-answer", handler.TimingAnswer)

	ec.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	ec.Logger.Fatal(ec.Start(":1323"))
}

func generateRandomSecret() string {
	b := make([]byte, 12)
	rand.Read(b)
	//randomSecretInt, _ := rand.Int(rand.Reader, new(big.Int).SetInt64(1000000000))
	randomSecret := fmt.Sprintf("%x", b)
	return randomSecret
}

func createOpenSqlite() {
	f, err := os.OpenFile("sqlite-database.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("Fatal error, unable to create sqlite file: %w \n", err))
	}
	f.Close()
	fmt.Println("sqlite db created")
}

func rotateSecrets(db *sql.DB) {
	var secret string
	var id int
	var deleteID int

	insertStudentSQL := `INSERT INTO secrets(secret) VALUES (?)`
	statement, perr := db.Prepare(insertStudentSQL)
	if perr != nil {
		fmt.Println("Fatal error, unable to insert to sqlite:", perr)
	}
	_, ierr := statement.Exec(generateRandomSecret())
	if ierr != nil {
		fmt.Println("Fatal error, unable to insert to sqlite:", ierr)
	}

	row, err := db.Query("SELECT id, secret FROM secrets ORDER BY id DESC")
	if err != nil {
		fmt.Println("error selecting from db")
	}
	defer row.Close()
	var rowCount int
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&id, &secret)
		fmt.Println("id & secret fetched from table")
		rowCount++
		if rowCount > 2 {
			if id > deleteID {
				deleteID = id
			}
		}
	}

	if deleteID > 0 {
		query := "DELETE FROM secrets WHERE id <= ?"
		_, r := db.Exec(query, deleteID)
		if r != nil {
			fmt.Println("error deleting:", r)
		}
	}
}

func hasSecretsTable(db *sql.DB) bool {
	row, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets';")
	if err != nil {
		return false
	}
	defer row.Close()
	var name string
	for row.Next() {
		row.Scan(&name)
	}

	if len(name) > 0 {
		return true
	}
	return false
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE table secrets ("id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,"secret" TEXT);`

	fmt.Println("Create secrets table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	fmt.Println("secrets table created")

	// insert a random secret
	insertStudentSQL := `INSERT INTO secrets(secret) VALUES (?)`
	statement, perr := db.Prepare(insertStudentSQL)
	if perr != nil {
		panic(fmt.Errorf("Fatal error, unable to insert to sqlite: %w \n", perr))
	}
	_, ierr := statement.Exec(generateRandomSecret())
	if ierr != nil {
		panic(fmt.Errorf("Fatal error, unable to insert to sqlite: %w \n", ierr))
	}
}

func insecureSecretGenerator() string {
	// this code tries to simulate issue outlined in https://github.com/YMFE/yapi/issues/2117

	// all values of base36
	//'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	//'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	//'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	//'U', 'V', 'W', 'X', 'Y', 'Z'
	//
	mrand.Seed(time.Now().UnixNano())
	randFloat := fmt.Sprintf("%v", mrand.Float64())[2:]
	// fmt.Println(randFloat)

	i, _ := strconv.Atoi(randFloat)

	//output is usually 11 bytes, we make it slightly easier with 10 bytes
	return base36.Encode(uint64(i))[0:10]
}
