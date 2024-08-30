package config

import (
	"fmt"
	"log"
	"os"
	"v-games-ip-ph2-ftgo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// err = godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connect to DB")
	}
	log.Println("Connected to DB")
	err = DB.AutoMigrate(models.User{},
		models.Game{},
		models.Review{},
		models.CouponCode{},
		models.DLC{},
		models.Cart{},
		models.CartItem{},
		models.Payment{},
		models.Rental{},
		models.RentalItem{})
	if err != nil {
		log.Fatal(err)
	}
}
