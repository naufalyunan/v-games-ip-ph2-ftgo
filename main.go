package main

import (
	"os"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitDB()

	//init echo
	e := echo.New()
	u := e.Group("/users")
	u.POST("/register", handlers.Register)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
