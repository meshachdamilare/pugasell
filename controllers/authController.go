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
var UserCollection *mongo.Collection = db.OpenCollection(db.Client, "user")

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"err": "error checking for email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exits"})
			return
		}

		// first registered user is an admin
		isFirstAccount, _ := UserCollection.CountDocuments(ctx, bson.M{})
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

		insertNum, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}

		// Generate token and attch it to cookie
		util.GenerateToken(c, user.Email, user.User_id, user.Role)

		// json response to user
		jsonResponse := make(map[string]interface{})
		jsonResponse["ID"] = insertNum
		jsonResponse["sucess"] = "user created"
		c.JSON(http.StatusOK, jsonResponse)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.LoginModel
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validateErr := validate.Struct(user); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No user found with the email"})
			return
		}
		passwordIsValid, msg := util.VerifyPassword(user.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		util.GenerateToken(c, foundUser.Email, foundUser.User_id, foundUser.Role)
		c.JSON(http.StatusOK, gin.H{"success": "you're logged in."})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "user logged out"})
	}
}
