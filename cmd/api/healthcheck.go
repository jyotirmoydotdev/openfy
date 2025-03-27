package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) healthCheckHandler(ctx *gin.Context) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(ctx.Writer, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(ctx, err)
	}
}
