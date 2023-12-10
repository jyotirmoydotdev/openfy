package models

type User struct {
	ID                string            `db:"id" gorm:"primaryKey"`
	Username          string            `db:"username" gorm:"index"`
	Password          string            `db:"password,omitempty" gorm:"uniqueIndex"`
	Email             string            `db:"email"`
	Phone             int               `db:"phone"`
	Age               int               `db:"age"`
	DeliveryAddresses []DeliveryAddress `db:"deliveryaddress" gorm:"foreignKey:UserID"`
}

type DeliveryAddress struct {
	UserID    string `db:"id" gorm:"primaryKey"`
	Country   string `db:"country"`
	Address   string `db:"address"`
	Apartment string `db:"apartment"`
	City      string `db:"city"`
	State     string `db:"state"`
	PinCode   int    `db:"pincode"`
}
