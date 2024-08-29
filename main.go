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

	i := e.Group("/items")
	i.Use(middlewares.IsAuthenticated("user"))
	i.POST("", handlers.CreateCartItem)

	c := e.Group("/carts")
	//both roles can view the carts
	c.Use(middlewares.IsAuthenticated("both"))
	c.GET("", handlers.GetCarts)

	co := e.Group("/coupons")
	co.Use(middlewares.IsAuthenticated("both"))
	co.GET("", handlers.GetCoupons)
	co.GET("/referral", handlers.GetReferralCode)

	p := e.Group("/payments")
	p.Use(middlewares.IsAuthenticated("user"))
	p.POST("/create", handlers.CreatePayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
