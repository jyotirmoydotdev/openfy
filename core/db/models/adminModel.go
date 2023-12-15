package models

type Admin struct {
	ID           string   `gorm:"column:id"`
	Username     string   `gorm:"column:username"`
	Password     string   `gorm:"column:password"`
	Email        string   `gorm:"column:email"`
	Name         string   `gorm:"column:name"`
	FirstName    string   `gorm:"column:firstName"`
	LastName     string   `gorm:"column:lastName"`
	AccountOwner bool     `gorm:"column:accountOwner"`
	Locale       string   `gorm:"column:locale"`
	Permission   []string `gorm:"column:permissions"`
}

var Admins []Admin
var AdminSecrets = make(map[string]string)
var AdminIDCounter int
