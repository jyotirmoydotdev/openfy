package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/auth"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/signup", auth.RegisterUser)
	router.POST("/login", auth.LoginUser)
	router.POST("/admin/signup", auth.RegisterAdmin)
	router.POST("/admin/login", auth.LoginAdmin)

	router.GET("/products", GetAllProducts)
	router.GET("/products/:id", GetProduct)
	authGroup := router.Group("/admin", auth.AuthenticateMiddleware())
	{
		authGroup.POST("/products/new", Create)
		authGroup.PUT("/products/:id", Update)
		authGroup.DELETE("/products/:id", Delete)
	}

	return router
}
