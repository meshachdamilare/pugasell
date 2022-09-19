package routes

import (
	controller "github.com/Christomesh/pugasell/controllers"
	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router *gin.Engine) {
	review := router.Group("/api/v1/reviews")
	{
		review.GET("/:review_id", controller.GetSingleReview())
		review.GET("/", controller.GetAllReviews())

		// ONly logged-in users can access the routes below
		review.Use(middleware.AuthenticateUser())

		review.POST("/", controller.CreateReview())
		review.PATCH("/:review_id", controller.UpdateReview())
		review.DELETE("/:review_id", controller.DeleteReview())

	}
}
