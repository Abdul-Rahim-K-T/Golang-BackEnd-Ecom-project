package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/raheem/Ecom/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {

	DB = connectDB()
	return DB
}

func connectDB() *gorm.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nConnected to DATABASE: ", DB.Name())
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Category_Offer{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.Image{})

	DB.AutoMigrate(&models.Brand{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Cart{})
	DB.AutoMigrate(&models.Wishlist{})
	DB.AutoMigrate(&models.Payment{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItem{})
	DB.AutoMigrate(&models.RazorPay{})
	DB.AutoMigrate(&models.Coupon{})

	return DB
}
