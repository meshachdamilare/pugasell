package routes

import (
	controller "github.com/Christomesh/pugasell/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/util")
	{
		auth.POST("/register", controller.Signup())
		auth.POST("/login", controller.Login())
		auth.GET("/logout", controller.Logout())
	}

}
