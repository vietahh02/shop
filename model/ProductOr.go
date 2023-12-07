package model

type ProductOr struct {
	Id string `json:"id"`
	Order
	Product
	Color
	Amount string
}
