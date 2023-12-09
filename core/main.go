package main

import (
	database "github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/internal/web"
)

func main() {
	router := web.SetupRouter()
	database.CreateDatabase()
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
