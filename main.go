package main

import (
	"fmt"
	"net/http"

	"weak-jwt/handler"

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

	ec := echo.New()
	ec.HideBanner = true

	ec.POST("/api/expired", handler.ExpiredLogin)
	ec.POST("/api/expired-answer", handler.ExpiredLoginAnswer)

	ec.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	ec.Logger.Fatal(ec.Start(":1323"))
}
