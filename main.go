package main

import (
	"os"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/handlers"
	"v-games-ip-ph2-ftgo/middlewares"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitDB()

	//init echo
	e := echo.New()
	u := e.Group("/users")
	u.POST("/register", handlers.Register)
	u.POST("/login", handlers.Login)

	g := e.Group("/games")
	g.GET("", handlers.GetGames)
	g.GET("/:id", handlers.GetGameByID)
	g.Use(middlewares.IsAuthenticated("admin"))
	g.POST("", handlers.CreateGame)
	g.PUT("/:id", handlers.UpdateGameStock)

	r := e.Group("/reviews")
	r.GET("", handlers.GetReviews)
	r.Use(middlewares.IsAuthenticated("user"))
	r.POST("", handlers.CreateReview)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
