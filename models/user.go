package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name string             `json:"first_name" validate:"required,min=2,max=50"`
	Last_name  string             `json:"last_name" validate:"required,min=2,max=50"`
	Email      string             `json:"email" validate:"required,email"`
	Avatar     string             `json:"avatar"`
	Password   string             `json:"password" validate:"required,min=6"`
	Role       string             `json:"role" default:"USER" validate:"omitempty,eq=ADMIN|eq=USER"`
	User_id    string             `json:"user_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}

type UserResponseModel struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name string             `json:"first_name" validate:"required,min=2,max=50"`
	Last_name  string             `json:"last_name" validate:"required,min=2,max=50"`
	Email      string             `json:"email" validate:"required,email"`
	Avatar     string             `json:"avatar"`
	Role       string             `json:"role" default:"USER" validate:"omitempty,eq=ADMIN|eq=USER"`
	User_id    string             `json:"user_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}

type LoginModel struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserModel struct {
	First_name string    `json:"first_name" validate:"required,min=2,max=50"`
	Last_name  string    `json:"last_name" validate:"required,min=2,max=50"`
	Email      string    `json:"email" validate:"required,email"`
	Avatar     string    `json:"avatar,omitempty"`
	Updated_at time.Time `json:"updated_at"`
}

// type ChangePasswordModel struct {
// 	OldPassword string `json:"old_password" validate:"required,min=6"`
// 	NewPassword string `json:"new_password" validate:"required,min=6"`
// }
