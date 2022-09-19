package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Christomesh/pugasell/db"
	"github.com/Christomesh/pugasell/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProductCollection *mongo.Collection = db.OpenCollection(db.Client, "product")

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(product); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		userId := c.GetString("userId")

		product.ID = primitive.NewObjectID()
		product.User_id = product.ID.Hex()
		product.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.User_id = userId

		insertNum, insertErr := ProductCollection.InsertOne(ctx, product)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ID": insertNum, "message": "success", "response": product})
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetSingleProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
