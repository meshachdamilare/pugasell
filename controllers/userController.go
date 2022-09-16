package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	middleware "github.com/Christomesh/pugasell/middleware"
	"github.com/Christomesh/pugasell/models"
	util "github.com/Christomesh/pugasell/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		//get userId, email, role  of logged-in user via details stored in the cookie
		// These can be used to get more info about the user from the database
		// And display them if neccessary

		user["userId"] = c.GetString("userId")
		user["email"] = c.GetString("email")
		user["role"] = c.GetString("role")
		c.JSON(http.StatusOK, gin.H{"success": user})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userUpdate models.UpdateUserModel
		var foundUser models.User

		if err := c.BindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validateErr := validate.Struct(userUpdate); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		//get userId of logged in user via details stored in the cookie
		userId := c.GetString("userId")
		role := c.GetString("role")

		var updateObj primitive.D

		if userUpdate.First_name != "" {
			updateObj = append(updateObj, bson.E{Key: "first_name", Value: userUpdate.First_name})

		}
		if userUpdate.Last_name != "" {
			updateObj = append(updateObj, bson.E{Key: "last_name", Value: userUpdate.Last_name})

		}
		if userUpdate.Email != "" {
			updateObj = append(updateObj, bson.E{Key: "email", Value: userUpdate.Email})

		}

		err := Usercollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		if err != nil {
			msg := "message:User was not found"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "avatar", Value: userUpdate.Avatar})
		userUpdate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: userUpdate.Updated_at})

		uspsert := true
		filter := bson.M{"user_id": userId}

		opt := options.UpdateOptions{
			Upsert: &uspsert,
		}

		result, err := Usercollection.UpdateOne(
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
		util.GenerateToken(c, userUpdate.Email, userId, role)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated", "result": result.ModifiedCount})

	}
}

// func UpdateUserPassword() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var changePassword models.ChangePasswordModel
// 		var foundUser models.User

// 		if err := c.BindJSON(&changePassword); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if validateErr := validate.Struct(changePassword); validateErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
// 			return
// 		}

// 		userId := c.GetString("userId")

// 		var updateObj primitive.D

// 		err := Usercollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
// 		if err != nil {
// 			msg := "message:User was not found"
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}
// 		isPasswordCorrect, msg := util.VerifyPassword(changePassword.OldPassword, foundUser.Password)
// 		if !isPasswordCorrect {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
// 		}

// 		if changePassword.NewPassword != "" {
// 			newPassword := util.HashPassword(changePassword.NewPassword)
// 			updateObj = append(updateObj, bson.E{Key: "password", Value: newPassword})

// 		}

// 		uspsert := true
// 		filter := bson.M{"user_id": userId}

// 		opt := options.UpdateOptions{
// 			Upsert: &uspsert,
// 		}

// 		result, err := Usercollection.UpdateOne(
// 			ctx,
// 			filter,
// 			bson.D{
// 				{Key: "$set", Value: updateObj},
// 			},
// 			&opt,
// 		)
// 		if err != nil {
// 			msg := "failed to change password"
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Password successfully changed", "result": result.ModifiedCount})
// 	}
// }
