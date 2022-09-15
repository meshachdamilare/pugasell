package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// limit, err := strconv.Atoi(c.Query("limit"))
		// if err != nil || limit < 1 {
		// 	limit = 2
		// }

		// page, err := strconv.Atoi(c.Query("page"))
		// if err != nil || page < 1 {
		// 	page = 1
		// }
		// skip := (page - 1) * limit
		length, err := Usercollection.CountDocuments(ctx, bson.M{"role": "USER"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		// to exclude the password field and avatar field from the result obtained and response object
		opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}, {Key: "avatar", Value: 0}})

		cursor, err := Usercollection.Find(ctx, bson.M{"role": "USER"}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		var allusers []bson.M
		if err := cursor.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{"users": allusers, "counts": length})
	}
}

func GetSingleUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ShowCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateUserPassword() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
