package model

type Info struct {
	Id          string `json:"id"`
	AccID       string `json:"acc_id"`
	UserName    string `json:"user_name"`
	FullName    string `json:"full_name"`
	Image       string `json:"image"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	BirthDay    string `json:"birth_day"`
	Gender      int    `json:"gender"`
	Address     string `json:"address"`
	Saddress    string `json:"saddress"`
}
