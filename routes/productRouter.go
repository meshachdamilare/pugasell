package routes

import (
	controller "github.com/Christomesh/pugasell/controllers"
	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	product := router.Group("/api/v1/products")
	{
		product.GET("/", controller.GetAllProducts())
		product.POST("/", middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"), controller.CreateProduct())
		product.POST("/uploadImage", middleware.AuthenticateUser(), controller.UploadImage())
		product.GET("/:product_id", middleware.AuthenticateUser(), controller.GetSingleProduct())
		product.PATCH("/:product_id", middleware.AuthenticateUser(), controller.UpdateProduct())
		product.DELETE("/:product_id", middleware.AuthenticateUser(), controller.DeleteProduct())
	}
}
