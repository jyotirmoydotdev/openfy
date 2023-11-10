package main

import "github.com/jyotirmoydotdev/openfy/web"

func main() {
	router := web.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
