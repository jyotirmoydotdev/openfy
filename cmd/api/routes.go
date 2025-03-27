package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	gin.DisableConsoleColor()
	router := gin.Default()

	router.HandleMethodNotAllowed = true

	router.NoMethod(app.methodNotAllowedResponse)
	router.NoRoute(app.notFoundResponse)

	router.GET("/v1/health", app.healthCheckHandler)

	return router
}
