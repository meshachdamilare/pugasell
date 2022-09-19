package main

import (
	"net/http"
	"os"

	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/Christomesh/pugasell/routes"
	"github.com/gin-gonic/gin"
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

	router.GET("/api/v1", middleware.AuthenticateUser(), middleware.AuthorizedPermissions("ADMIN"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to PUGASELL"})
	})

	router.Run(":" + port)

}
