package model

type Comment struct {
	Id      string `json:"id"`
	Info    `json:"info`
	Product `json:"product"`
	Color   `json:"color"`
	Star    string `json:"star"`
	Com     string `json:"com"`
	Date    string `json:"date"`
}
