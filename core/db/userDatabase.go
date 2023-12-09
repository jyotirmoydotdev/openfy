package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

func CreateUserDatabase() {
	// Open SQLite database (create if not exists)
	db, err := sqlx.Open("sqlite3", "./Database/user.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create User table
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		username TEXT,
		password TEXT,
		email TEXT,
		phone INTEGER,
		age INTEGER
	);`

	_, err = db.Exec(createUserTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Create Delivery table
	createDeliveryTableSQL := `
	CREATE TABLE IF NOT EXISTS delivery (
		user_id TEXT,
		country TEXT,
		address TEXT,
		apartment TEXT,
		city TEXT,
		state TEXT,
		pincode INTEGER,
		FOREIGN KEY (user_id) REFERENCES user(id)
	);`

	_, err = db.Exec(createDeliveryTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
