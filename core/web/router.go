package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/auth"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public Route
	router.POST("/signup", auth.RegisterUser)
	router.POST("/login", auth.LoginUser)
	router.GET("/products", GetAllActiveProducts)

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
		// user.GET("/orders",GetAllOrders)
		// user.GET("/orders:id", GetOrder)
		// user.GET("/profile", profile)
		// user.PUT("/profile", UpdateProfile)
	}

	// Admin Route
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
