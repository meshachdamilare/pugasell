package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/pugasell/db"
	"github.com/meshachdamilare/pugasell/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		length, err := ProductCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		cursor, err := ProductCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		var productResponse []models.Product
		if err := cursor.All(ctx, &productResponse); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"products": productResponse, "counts": length})
	}
}

func GetSingleProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))
		var productResponse models.ProductResponse

		err := ProductCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&productResponse.Product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		errReview := ReviewCollection.FindOne(ctx, bson.M{"product_id": productId}).Decode(&productResponse.Review)
		if errReview != nil {
			if errReview == mongo.ErrNoDocuments {
				c.JSON(http.StatusOK, productResponse)
				return
			}
			return
		}

		c.JSON(http.StatusOK, productResponse)
	}
}

func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productUpdate models.UpdateProduct
		var foundProduct models.Product
		productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))

		if err := c.BindJSON(&productUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validateErr := validate.Struct(productUpdate); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		var updateObj primitive.D

		if productUpdate.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: productUpdate.Name})

		}
		if productUpdate.Price != 0 {
			updateObj = append(updateObj, bson.E{Key: "price", Value: productUpdate.Price})

		}
		if productUpdate.Description != "" {
			updateObj = append(updateObj, bson.E{Key: "description", Value: productUpdate.Description})

		}
		if productUpdate.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: productUpdate.Category})

		}
		if productUpdate.Company != "" {
			updateObj = append(updateObj, bson.E{Key: "company", Value: productUpdate.Company})

		}

		err := UserCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&foundProduct)
		if err != nil {
			msg := fmt.Sprintf("No product found with the id %s", productId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "image", Value: productUpdate.Image})
		updateObj = append(updateObj, bson.E{Key: "colors", Value: productUpdate.Colors})
		updateObj = append(updateObj, bson.E{Key: "featured", Value: productUpdate.Featured})
		updateObj = append(updateObj, bson.E{Key: "freeShipping", Value: productUpdate.FreeShipping})
		updateObj = append(updateObj, bson.E{Key: "inventory", Value: productUpdate.Inventory})

		productUpdate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: productUpdate.Updated_at})

		uspsert := true
		filter := bson.M{"_d": productId}

		opt := options.UpdateOptions{
			Upsert: &uspsert,
		}

		result, err := ProductCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "user item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product successfully updated", "result": result.ModifiedCount})

	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))

		result, err := ProductCollection.DeleteOne(ctx, bson.M{"_id": productId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "product deleted", "count": result.DeletedCount})
	}
}
