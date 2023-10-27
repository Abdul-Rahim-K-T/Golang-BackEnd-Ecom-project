package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	_ "github.com/raheem/Ecom/docs"
	"github.com/raheem/Ecom/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ecom API
//@version 1.0
//@discription Ecom API in go using Gin frame work

// @host     localhost:8080
// @BasePath/
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	database.InitDB()

	router := gin.Default()

	router.Static("/public", "./public")

	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	//add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + port)
}
