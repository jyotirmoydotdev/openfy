package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/auth"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/signup", auth.RegisterUser)
	router.POST("/login", auth.LoginUser)
	router.POST("/admin/signup", hashAdmin(), auth.RegisterAdmin)
	router.POST("/admin/login", auth.LoginAdmin)

	user := router.Group("/api/v1", auth.AuthenticateUserMiddleware())
	{
		user.GET("/products", GetAllProducts)
	}

	admin := router.Group("/admin", auth.AuthenticateMiddleware())
	{
		admin.GET("/products/:id", GetProduct)
		admin.GET("/products", GetAllProducts)
		admin.POST("/products/new", Create)
		admin.PUT("/products/:id", Update)
		admin.DELETE("/products/:id", Delete)
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
