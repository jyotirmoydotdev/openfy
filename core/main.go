package main

import (
	"github.com/jyotirmoydotdev/openfy/internal/web"
)

func main() {
	router := web.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
