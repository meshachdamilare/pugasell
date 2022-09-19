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
		product.GET("/:product_id", controller.GetSingleProduct())

		// ONly admin can access this resource
		product.Use(middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"))

		product.POST("/", controller.CreateProduct())
		product.POST("/uploadImage", controller.UploadImage())
		product.PATCH("/:product_id", controller.UpdateProduct())
		product.DELETE("/:product_id", controller.DeleteProduct())
	}
}
