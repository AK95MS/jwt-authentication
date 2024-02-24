package routes

import (
	controller "go_jwt_authentication/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("user/signup", controller.Sighup())
	incomingRoutes.POST("users/login", controller.Login())
}
