package routes

import (
	controller "go_jwt_authentication/controller"
	"go_jwt_authentication/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("users/:user_id", controller.GetUser())
}
