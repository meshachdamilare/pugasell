package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Christomesh/pugasell/db"
	"github.com/Christomesh/pugasell/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ReviewCollection *mongo.Collection = db.OpenCollection(db.Client, "review")

func CreateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		var review models.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(review); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		productId := review.Product_id
		primitiveProductId, _ := primitive.ObjectIDFromHex(productId) // convert it from string to ObjectId() used in mongodb
		userId := c.GetString("userId")                               // get the user_id of the logged-in user via jwt_token payload/claims

		// check if the product_id passed exist in the Product Collection
		isValidProduct := ProductCollection.FindOne(ctx, bson.M{"_id": primitiveProductId})

		if isValidProduct.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": isValidProduct.Err().Error()})
			return
		}

		// to ensure a user only submit a review per product
		filter := bson.M{"user_id": userId, "product_id": review.Product_id}

		alreadySubmitterReview := ReviewCollection.FindOne(ctx, filter)
		if alreadySubmitterReview.Err() == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already submitted review for this product"})
			return
		}

		review.ID = primitive.NewObjectID()
		review.User_id = userId
		review.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		review.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		insertNum, insertErr := ReviewCollection.InsertOne(ctx, review)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "review was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ID": insertNum, "message": "success", "response": review})

	}
}

func GetAllReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		length, err := ReviewCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		cursor, err := ReviewCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		var reviewResponse []models.Review
		if err := cursor.All(ctx, &reviewResponse); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"products": reviewResponse, "counts": length})
	}
}

func GetSingleReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		reviewId, _ := primitive.ObjectIDFromHex(c.Param("review_id"))
		var reviewResponse models.Review

		err := ReviewCollection.FindOne(ctx, bson.M{"_id": reviewId}).Decode(&reviewResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, reviewResponse)
	}
}

func UpdateReview() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
