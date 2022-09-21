package controllers

import (
	"context"
	"fmt"
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

var OrderCollection *mongo.Collection = db.OpenCollection(db.Client, "order")

func fakeStripeAPI(amount float64, currency string) (string, float64) {
	var client_secret = "someRandomValue"
	return client_secret, amount
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		var order models.OrderItems
		var singleOrderItem models.SingleOrderItem
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(order); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		userId := c.GetString("userId")

		if order.Tax == 0.0 && order.ShippingFee == 0.0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "please provide tax and shipping fee"})
			return
		}

		var orderItems []models.SingleOrderItem
		var subtotal float64 = 0
		var total float64 = 0

		for _, item := range order.Items {
			primitiveProductId, _ := primitive.ObjectIDFromHex(item.Product_id)
			err := ProductCollection.FindOne(ctx, bson.M{"_id": primitiveProductId}).Decode(&singleOrderItem)
			if err != nil {
				msg := fmt.Sprintf("No product with id %s", item.Product_id)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			singleOrderItem.Amount = item.Amount
			singleOrderItem.Name = item.Name
			singleOrderItem.Price = item.Price
			singleOrderItem.Image = item.Image
			singleOrderItem.Product_id = item.Product_id

			orderItems = append(orderItems, singleOrderItem)

			subtotal += item.Amount * item.Price

		}
		total = order.Tax + order.ShippingFee + subtotal

		clientSecret, amount := fakeStripeAPI(total, "usd")

		order.ClientSecret = clientSecret
		order.Total = amount
		order.SubTotal = subtotal
		order.Items = orderItems
		order.User_id = userId
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		insertNum, insertErr := OrderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ID": insertNum, "message": "success", "response": order})
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		length, err := OrderCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		cursor, err := OrderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		var orderResponse []models.OrderItems
		if err := cursor.All(ctx, &orderResponse); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"products": orderResponse, "counts": length})
	}
}

func GetSingleOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetCurrentUserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdatedOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
