package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//	{
//	    "name": "accent chair",
//	    "price": 25999,
//	    "image": "https://dl.airtable.com/.attachmentThumbnails/e8bc3791196535af65f40e36993b9e1f/438bd160",
//	    "colors": ["#ff0000", "#00ff00", "#0000ff"],
//	    "company": "marcos",
//	    "description": "Cloud bread VHS hell of banjo bicycle rights jianbing umami mumblecore etsy 8-bit pok pok +1 wolf. Vexillologist yr dreamcatcher waistcoat, authentic chillwave trust fund. Viral typewriter fingerstache pinterest pork belly narwhal. Schlitz venmo everyday carry kitsch pitchfork chillwave iPhone taiyaki trust fund hashtag kinfolk microdosing gochujang live-edge",
//	    "category": "office"
//	}
type Product struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `json:"name" validate:"required,min=2,max=100"`
	Price           float64            `json:"price" default:"0" validate:"required"`
	Description     string             `json:"description" validate:"required,max=1000"`
	Image           string             `json:"image" default:"/uploads/example.jpeg"`
	Category        string             `json:"category" validate:"required,eq=kitchen|eq=office|eq=bedroom"`
	Company         string             `json:"company" validate:"required,eq=ikea|eq=liddy|eq=marcos"`
	Colors          []string           `json:"colors" default:"[#2222]" validate:"required"`
	Featured        bool               `json:"featured" default:"false"`
	FreeShipping    bool               `json:"freeShipping" default:"false"`
	Inventory       float64            `json:"inventory,omitempty" default:"15"`
	AverageRating   float64            `json:"averageRating" default:"0"`
	NumberOfReviews float64            `json:"numberOfReviews" default:"0"`
	User_id         string             `json:"user_id"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
	//Reviews         []Review           `json:"reviews" bson:"-"` // we dont strore reviews in the db... this will be used as response to the user
}

type UpdateProduct struct {
	Name         string    `json:"name" validate:"required,min=2,max=100"`
	Price        float64   `json:"price" default:"0" validate:"required"`
	Description  string    `json:"description" validate:"required,max=1000"`
	Image        string    `json:"image" default:"/uploads/example.jpeg"`
	Category     string    `json:"category" validate:"required,eq=kitchen|eq=office|eq=bedroom"`
	Company      string    `json:"company" validate:"required,eq=ikea|eq=liddy|eq=marcos"`
	Colors       []string  `json:"colors" default:"[#2222]" validate:"required"`
	Featured     bool      `json:"featured" default:"false"`
	FreeShipping bool      `json:"freeShipping" default:"false"`
	Inventory    float64   `json:"inventory,omitempty" default:"15"`
	Updated_at   time.Time `json:"updated_at"`
}

type ProductResponse struct {
	Product Product
	Review  []Review `json:"reviews"`
}
