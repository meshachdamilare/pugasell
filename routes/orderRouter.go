package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/meshachdamilare/pugasell/controllers"
	middleware "github.com/meshachdamilare/pugasell/middleware"
)

func OrderRoutes(router *gin.Engine) {
	order := router.Group("/api/v1/orders")
	{
		order.GET("/", middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"), controller.GetAllOrders())

		order.POST("/", middleware.AuthenticateUser(), controller.CreateOrder())
		order.GET("/showAllMyOrders", middleware.AuthenticateUser(), controller.GetCurrentUserOrders())
		order.GET("/:order_id", middleware.AuthenticateUser(), controller.GetSingleOrder())
		order.PATCH("/:order_id", middleware.AuthenticateUser(), controller.UpdatedOrder())

	}
}
