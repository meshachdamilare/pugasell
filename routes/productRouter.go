package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/meshachdamilare/pugasell/controllers"
	middleware "github.com/meshachdamilare/pugasell/middleware"
)

func ProductRoutes(router *gin.Engine) {
	product := router.Group("/api/v1/products")
	{
		product.GET("/", controller.GetAllProducts())
		product.GET("/:product_id", controller.GetSingleProduct())
		//product.GET("/:product_id/", controller.GetSingleProductReviews())

		// ONly admin can access this resource
		product.Use(middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"))

		product.POST("/", controller.CreateProduct())
		product.POST("/uploadImage", controller.UploadImage())
		product.PATCH("/:product_id", controller.UpdateProduct())
		product.DELETE("/:product_id", controller.DeleteProduct())
	}
}
