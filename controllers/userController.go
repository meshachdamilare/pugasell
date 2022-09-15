package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/Christomesh/pugasell/models"
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

		// to exclude the password field from the result obtained and response object
		opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

		cursor, err := Usercollection.Find(ctx, bson.M{"role": "USER"}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		var userResponse []models.UserResponseModel
		if err := cursor.All(ctx, &userResponse); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{"users": userResponse, "counts": length})
	}
}

func GetSingleUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var userResponse models.UserResponseModel

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := middleware.CheckPermission(c, userId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not Authorized to access this route"})
			return
		}

		opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
		err := Usercollection.FindOne(ctx, bson.M{"user_id": userId}, opts).Decode(&userResponse)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": userResponse})
	}
}

func ShowCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = make(map[string]string)

		user["userId"] = c.GetString("userId")
		user["email"] = c.GetString("email")
		user["role"] = c.GetString("role")
		c.JSON(http.StatusOK, gin.H{"success": user})
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
