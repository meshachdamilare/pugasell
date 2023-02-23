package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/meshachdamilare/pugasell/controllers"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/util")
	{
		auth.POST("/register", controller.Signup())
		auth.POST("/login", controller.Login())
		auth.GET("/logout", controller.Logout())
	}

}
