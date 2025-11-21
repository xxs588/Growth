package routes

import (
	"mygo/controller"

	"github.com/gin-gonic/gin"
)

func InintUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/send_code", controller.SendCode)
		userGroup.POST("/register", controller.Register)
		userGroup.POST("/login", controller.Login)
	}
}
