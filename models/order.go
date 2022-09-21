package models

import "time"

type SingleOrderItem struct {
	Name       string  `json:"name" validate:"required"`
	Image      string  `json:"image" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Product_id string  `json:"product_id" validate:"required"`
}

type OrderItems struct {
	Tax             float64           `json:"tax"`
	ShippingFee     float64           `json:"shipping_fee"`
	SubTotal        float64           `json:"subtotal"`
	Total           float64           `json:"total"`
	Items           []SingleOrderItem `json:"items"`
	Status          string            `json:"status" default:"pending" validate:"omitempty,eq=pending|eq=failed|eq=paid|eq=delivered|eq=canceled"`
	User_id         string            `json:"user_id"`
	ClientSecret    string            `json:"client_secret"`
	PaymentIntentId string            `json:"payment_id"`
	Created_at      time.Time         `json:"created_at"`
	Updated_at      time.Time         `json:"updated_at"`
}

// type OrderResponse struct {
// 	Tax          *float64          `json:"tax" `
// 	ShippingFee  *float64          `json:"shipping_fee"`
// 	SubTotal     float64           `json:"subtotal"`
// 	Total        float64           `json:"total"`
// 	Items        []SingleOrderItem `json:"items"`
// 	Status       string            `json:"status"`
// 	User_id      string            `json:"user_id"`
// 	ClientSecret string            `json:"client_secret"`
// }
