package model

type Color struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	ColorUrl  string `json:"color_url"`
	Discount  string `json:"discount"`
}
