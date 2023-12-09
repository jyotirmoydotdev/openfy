package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/internal/auth"
	web "github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	requestProductHandlers := web.NewRequestProductHandlers()
	// Public Route
	router.POST("/signup", auth.RegisterUser)
	router.POST("/login", auth.LoginUser)
	// router.GET("/products", GetAllActiveProducts)

	router.POST("/admin/signup", hashAdmin(), auth.RegisterAdmin)
	router.POST("/admin/login", auth.LoginAdmin)

	// User route
	user := router.Group("/api", auth.AuthenticateUserMiddleware())
	{
		user.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Admin Route
	admin := router.Group("/admin", auth.AuthenticateMiddleware())
	{
		admin.GET("/products/:id", requestProductHandlers.GetProduct)
		admin.GET("/products", requestProductHandlers.GetAllProducts)
		admin.POST("/products/new", requestProductHandlers.Create)
		admin.PUT("/products/:id", requestProductHandlers.Update)
		admin.DELETE("/products/:id", requestProductHandlers.Delete)
	}
	return router
}

func hashAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashAdmin := auth.HashAdmin()
		if hashAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Contact Admin for signup",
			})
			ctx.Abort()
			return
		}
	}
}
