package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Christomesh/pugasell/db"
	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/Christomesh/pugasell/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			msg := fmt.Sprintf("no product found with id %s", productId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		reviewId, _ := primitive.ObjectIDFromHex(c.Param("review_id"))

		var reviewUpdate models.UpdateReview
		var foundReview models.Review

		if err := c.BindJSON(&reviewUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(reviewUpdate); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := ReviewCollection.FindOne(ctx, bson.M{"_id": reviewId}).Decode(&foundReview)
		if err != nil {
			msg := fmt.Sprintf("No product found with the id %s", reviewId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		userId := foundReview.User_id

		if err := middleware.CheckPermission(c, userId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not Authorized to access this route"})
			return
		}

		var updateObj primitive.D

		if reviewUpdate.Title != "" {
			updateObj = append(updateObj, bson.E{Key: "title", Value: reviewUpdate.Title})

		}
		if reviewUpdate.Comment != "" {
			updateObj = append(updateObj, bson.E{Key: "comment", Value: reviewUpdate.Comment})

		}

		updateObj = append(updateObj, bson.E{Key: "rating", Value: reviewUpdate.Rating})

		reviewUpdate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: reviewUpdate.Updated_at})

		uspsert := true
		filter := bson.M{"_id": reviewId}

		opt := options.UpdateOptions{
			Upsert: &uspsert,
		}

		result, err := ReviewCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "review item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Review successfully updated", "result": result.ModifiedCount})

	}
}

func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var foundReview models.Review
		reviewId, _ := primitive.ObjectIDFromHex(c.Param("review_id"))

		err := ReviewCollection.FindOne(ctx, bson.M{"_id": reviewId}).Decode(&foundReview)
		if err != nil {
			msg := fmt.Sprintf("No product found with the id %s", c.Param("review_id"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		userId := foundReview.User_id

		if err := middleware.CheckPermission(c, userId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not Authorized to access this route"})
			return
		}

		result, err := ReviewCollection.DeleteOne(ctx, bson.M{"_id": reviewId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "review deleted", "count": result.DeletedCount})
	}
}

func GetSingleProductReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productId := c.Param("product_id")
		var productReviews models.Review

		cursor, err := ReviewCollection.Find(ctx, bson.M{"product_id": productId})
		if err != nil {
			msg := fmt.Sprintf("No product with is %s", productId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if err = cursor.All(ctx, &productReviews); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, productReviews)
	}
}
