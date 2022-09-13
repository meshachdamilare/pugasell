package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Christomesh/pugasell/db"
	"github.com/Christomesh/pugasell/models"
	util "github.com/Christomesh/pugasell/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var Usercollection *mongo.Collection = db.OpenCollection(db.Client, "user")

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := Usercollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"err": "error checking for email"})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exits"})
		}

		// first registered user is an admin
		isFirstAccount, _ := Usercollection.CountDocuments(ctx, bson.M{})
		if isFirstAccount == 0 {
			user.Role = "ADMIN"
		} else {
			user.Role = "USER"
		}

		password := util.HashPassword(user.Password)
		user.Password = password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		insertNum, insertErr := Usercollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
		}

		c.JSON(http.StatusOK, insertNum)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
