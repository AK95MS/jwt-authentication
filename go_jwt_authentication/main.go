package main

import (
	"fmt"
	routes "go_jwt_authentication/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("welcome to JWT authentication by GO")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{"success": "Access Granted For Api-1"})
	})

	router.GET("API-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"suceess": "Access Granted For API-2"})
	})

	router.Run(":" + port)

}
