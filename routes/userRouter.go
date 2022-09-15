package routes

import (
	controller "github.com/Christomesh/pugasell/controllers"
	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/api/v1/users")
	{
		user.GET("/", middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"), controller.GetAllUsers())
		user.GET("/showme", middleware.AuthenticateUser(), controller.ShowCurrentUser())
		user.GET("/:id", middleware.AuthenticateUser(), controller.GetSingleUser())
		user.PATCH("/updateuser", middleware.AuthenticateUser(), controller.UpdateUser())
		user.PATCH("/updatepassword", middleware.AuthenticateUser(), controller.UpdateUserPassword())
	}
}
