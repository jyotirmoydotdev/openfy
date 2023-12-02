package database

type User struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	Password        string     `json:"password,omitempty"`
	Email           string     `json:"email"`
	Phone           int        `json:"phone"`
	Age             int        `json:"age"`
	DeliveryAddress []Delivery `json:"deliveryaddress"`
}

type Delivery struct {
	Country   string `json:"country"`
	Address   string `json:"address"`
	Apartment string `json:"apartment"`
	City      string `json:"city"`
	State     string `json:"statte"`
	PinCode   int    `json:"pincode"`
}

var Users []User
var UserSecrets = make(map[string]string)
