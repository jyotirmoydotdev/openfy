package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/jyotirmoydotdev/openfy/db/models"
	_ "github.com/mattn/go-sqlite3"
)

var Users []models.User
var UserSecrets = make(map[string]string)

func CreateUserDatabase() {
	// Open SQLite database (create if not exists)
	db, err := sqlx.Open("sqlite3", "./db/data/user.db")
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
