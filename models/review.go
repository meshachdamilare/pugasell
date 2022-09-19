package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID         primitive.ObjectID `bson:"_id"`
	Rating     float64            `json:"rating" validate:"required,min=1,max=5"`
	Title      string             `json:"title" validate:"required,max=50"`
	Comment    string             `json:"comment" validate:"required"`
	User_id    string             `json:"user_id"`
	Product_id string             `json:"product_id"`
}
