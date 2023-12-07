package model

type OrderBeta struct {
	Id   string
	OsId string
	Store
	ListProduct []Product
	Status
	Quantity int
	AllPrice float64
}
