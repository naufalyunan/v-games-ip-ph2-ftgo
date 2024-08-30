package main

import (
	"os"
	"v-games-ip-ph2-ftgo/config"
	_ "v-games-ip-ph2-ftgo/docs"
	"v-games-ip-ph2-ftgo/handlers"
	"v-games-ip-ph2-ftgo/middlewares"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Kedai Game API
// @version 1.0
// @description Ini adalah API untuk rental video game
// @termOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@kedaiGame.com
// @host localhost:8080
// @basePath /

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	config.InitDB()

	//init echo
	e := echo.New()

	//Register Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//users routes
	u := e.Group("/users")
	u.POST("/register", handlers.Register)
	u.POST("/login", handlers.Login)
	u.Use(middlewares.IsAuthenticated("user"))
	u.PUT("/depo", handlers.TopUp)

	//games routes
	g := e.Group("/games")
	g.GET("", handlers.GetGames)
	g.GET("/:id", handlers.GetGameByID)
	g.Use(middlewares.IsAuthenticated("admin"))
	g.POST("", handlers.CreateGame)
	g.PUT("/:id", handlers.UpdateGameStock)

	//reviews routes
	r := e.Group("/reviews")
	r.GET("", handlers.GetReviews)
	r.Use(middlewares.IsAuthenticated("user"))
	r.POST("", handlers.CreateReview)

	//items routes
	i := e.Group("/items")
	i.Use(middlewares.IsAuthenticated("user"))
	i.POST("", handlers.CreateCartItem)

	//carts routes
	c := e.Group("/carts")
	//both roles can view the carts
	c.Use(middlewares.IsAuthenticated("both"))
	c.GET("", handlers.GetCarts)

	//coupons routes
	co := e.Group("/coupons")
	co.Use(middlewares.IsAuthenticated("both"))
	co.GET("", handlers.GetCoupons)
	co.GET("/referral", handlers.GetReferralCode)

	//payments routes
	p := e.Group("/payments")
	p.Use(middlewares.IsAuthenticated("both"))
	p.GET("", handlers.GetPayments)
	p.Use(middlewares.IsAuthenticated("user"))
	p.POST("/create", handlers.CreatePayment)

	//pay payments routes
	pay := e.Group("/pay")
	pay.Use(middlewares.IsAuthenticated("admin"))
	pay.PUT("/:id", handlers.Pay)

	//rentals routes
	rt := e.Group("/rentals")
	rt.Use(middlewares.IsAuthenticated("both"))
	rt.GET("", handlers.GetRentals)

	//rental-items routes
	ri := e.Group("/rental-item")
	ri.Use(middlewares.IsAuthenticated("admin"))
	ri.PUT("/:id", handlers.UpdateRentalItem)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
