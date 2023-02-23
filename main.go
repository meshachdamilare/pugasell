package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	middleware "github.com/meshachdamilare/pugasell/middleware"
	"github.com/meshachdamilare/pugasell/routes"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.ProductRoutes(router)
	routes.ReviewRoutes(router)
	routes.OrderRoutes(router)

	router.GET("/api/v1", middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to PUGASELL"})
	})

	router.Run(":" + port)

}
