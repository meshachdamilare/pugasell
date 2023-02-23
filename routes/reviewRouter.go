package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/meshachdamilare/pugasell/controllers"
	middleware "github.com/meshachdamilare/pugasell/middleware"
)

func ReviewRoutes(router *gin.Engine) {
	review := router.Group("/api/v1/reviews")
	{
		review.GET("/:review_id", controller.GetSingleReview())
		review.GET("/", controller.GetAllReviews())

		// ONly logged-in users can access the route below
		review.Use(middleware.AuthenticateUser())

		review.POST("/", controller.CreateReview())
		review.PATCH("/:review_id", controller.UpdateReview())
		review.DELETE("/:review_id", controller.DeleteReview())

	}
}
