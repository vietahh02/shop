package model

type OrderDetailBeta struct {
	Id   string
	OsId string
	Store
	ListProduct []Product
	Status
	Quantity int
	AllPrice float64
	Info
	Delivery
	Payment
	DateOrder   string
	DateReceive string
}
