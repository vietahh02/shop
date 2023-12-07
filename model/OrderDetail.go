package model

type ProductOrder struct {
	ProductName string
	ColorName   string
	Amount      string
	Price       string
}

type OrderDetail struct {
	Id      string `json:"id"`
	OrderID string `json:"order_id"`
	Store
	ListPro []ProductOrder
	Status
}
