package model

type Cart struct {
	Id        string `json:"id"`
	CartID    string `json:"cart_id"`
	ProductID string `json:"product_id"`
	ColorID   string `json:"color_id"`
	Amount    int    `json:"amount"`
}
