package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/meshachdamilare/pugasell/controllers"
	middleware "github.com/meshachdamilare/pugasell/middleware"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/api/v1/users")
	{

		// Omly logged-in users can access the route below
		user.Use(middleware.AuthenticateUser())

		user.GET("/showMe", controller.ShowCurrentUser())
		user.GET("/:userId", controller.GetSingleUser())
		user.PATCH("/updateuser", controller.UpdateUser())
		user.PATCH("/updatepassword", controller.UpdateUserPassword())

		// Only Admin can access this route
		user.Use(middleware.AuthorizedPermissions("ADMIN"))
		user.GET("/", controller.GetAllUsers())

		// Another way to use the util
		// user.GET("/showMe",util.AuthenticateUser(), controller.ShowCurrentUser())
	}
}
