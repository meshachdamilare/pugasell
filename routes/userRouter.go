package routes

import (
	controller "github.com/Christomesh/pugasell/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/api/v1/user")
	{
		user.GET("/", controller.GetAllUsers())
		user.GET("/showme", controller.ShowCurrentUser())
		user.GET("/:id", controller.GetSingleUser())
		user.PATCH("/updateuser", controller.UpdateUser())
		user.PATCH("/updatepassword", controller.UpdateUserPassword())
	}
}
