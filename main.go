package main

import (
	"fmt"
	"net/http"

	"github.com/jasonwubz/weak-jwt/handlers"

	"github.com/labstack/echo/v4"
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

	e := echo.New()
	e.HideBanner = true

	e.POST("/api/Expired", handlers.ExpiredLogin)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
