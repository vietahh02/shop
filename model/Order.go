package model

type Order struct {
	Id           string `json:"id"`
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Saddress     string `json:"saddress"`
	Date         string `json:"date"`
	Status
	Delivery
	Payment
	OrderDetailID string
}
