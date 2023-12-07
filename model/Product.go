package model

type Product struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Price          float64  `json:"price"`
	Discount       float64  `json:"discount"`
	Star           string   `json:"star"`
	Image          string   `json:"image"`
	Date           string   `json:"date"`
	Sold           int      `json:"sold"`
	InStock        int      `json:"in_stock"`
	Description    string   `json:"description"`
	Sub_Image      []string `json:"sub_image"`
	Discount_Color string   `json:"discount_color"`
	Category       `json:"category"`
	Color
}
