package controllers

import (
	"github.com/Christomesh/pugasell/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var ReviewCollection *mongo.Collection = db.OpenCollection(db.Client, "review")

func CreateReview() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetAllReviews() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetSingleReview() gin.HandlerFunc {
	return func(c *gin.Context) {

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
