package main

import (
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/internal/web"
)

func main() {
	router := web.SetupRouter()
	err := db.InitializeDatabases()
	if err != nil {
		panic(err)
	}
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
