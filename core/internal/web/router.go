package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/internal/auth"
	"github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

func SetupRouter() *gin.Engine {
	err := db.InitializeDatabases()
	if err != nil {
		panic(err)
	}
	router := gin.Default()

	router.POST("/signup", auth.RegisterUser)
	router.POST("/login", auth.LoginUser)
	// router.GET("/products", GetAllActiveProducts)

	router.POST("/admin/signup", hashAdmin(), auth.SignupAdmin)
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

	admin := router.Group("/admin", auth.AuthenticateMiddleware())
	{
		admin.GET("/products/:id", handlers.GetProduct)
		admin.GET("/products", handlers.GetAllProducts)
		admin.POST("/products/new", handlers.Create)
		admin.PUT("/products/:id", handlers.Update)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
		admin.DELETE("/products/:id/:vid", handlers.DeleteProductVarient)
		// admin.POST("/auth-with-password", AuthWithPassword)
		// admin.POST("/request-password-reset", RequestPasswordReset)
		// admin.POST("/confirm-password-reset", ConfirmPasswordReset)
		// admin.POST("/auth-refresh", AuthRefresh)
		// admin.GET("", list) // get Admin List
		// admin.POST("", create) // create new admin
		// admin.GET("/:id", view) // view a admin detail
		// admin.PATCH("/:id", update) // update the admin details
		// admin.DELETE("/:id", delete) // delete the admin
	}
	return router
}

func hashAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashNoAdmin, err := auth.HashAdmin()
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Internal Server Error",
			})
			ctx.Abort()
			return
		}
		if !hashNoAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Contact Admin for signup",
			})
			ctx.Abort()
			return
		}
	}
}
