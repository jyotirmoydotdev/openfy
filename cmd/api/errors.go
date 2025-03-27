package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(ctx *gin.Context, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(ctx.Writer, status, env, nil)
	if err != nil {
		app.logError(ctx.Request, err)
		ctx.Writer.WriteHeader(500)
	}

}

func (app *application) serverErrorResponse(ctx *gin.Context, err error) {
	app.logError(ctx.Request, err)
	message := "the server encountered a problem an could not process your request"
	app.errorResponse(ctx, http.StatusInternalServerError, message)
}

func (app *application) methodNotAllowedResponse(ctx *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for the resource", ctx.Request.Method)
	app.logger.Println("Method Not Allowed:", ctx.Request.Method)
	app.errorResponse(ctx, http.StatusMethodNotAllowed, message)
}

func (app *application) notFoundResponse(ctx *gin.Context) {
	message := "the request resource could not be found"
	app.errorResponse(ctx, http.StatusNotFound, message)
}
