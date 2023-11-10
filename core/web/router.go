package web

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/products/new", Create)
	router.GET("/products", GetAllProducts)
	router.GET("/products/:id", GetProduct)
	router.PUT("/products/:id", Update)
	router.DELETE("/products/:id", Delete)

	return router
}
