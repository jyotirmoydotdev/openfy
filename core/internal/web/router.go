package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jyotirmoydotdev/openfy/database"
	"github.com/jyotirmoydotdev/openfy/internal/auth"
	"github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

func SetupRouter() *gin.Engine {
	err := db.InitializeDatabases()
	if err != nil {
		panic(err)
	}
	router := gin.Default()

	router.POST("/signup", auth.RegisterCustomer)
	router.POST("/login", auth.LoginCustomer)
	router.GET("/products", handlers.GetAllActiveProducts)

	router.POST("/staffMember/signup", hashStaffMember(), auth.SignupStaffMember)
	router.POST("/staffMember/login", auth.LoginStaffMember)

	// Customer route
	customer := router.Group("/customer", auth.AuthenticateCustomerMiddleware())
	{
		customer.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	staffMember := router.Group("/staffMember", auth.AuthenticateMiddleware())
	{
		// GET example.com/staffMember/product?id=x
		staffMember.GET("/product", handlers.GetProduct)

		staffMember.GET("/products", handlers.GetAllProducts)
		staffMember.POST("/product/new", handlers.Create)

		// PUT example.com/staffMember/product?id=x
		staffMember.PUT("/product", handlers.Update)

		// DELETE example.com/staffMember/product?id=x
		staffMember.DELETE("/product", handlers.DeleteProduct)

		// DELETE example.com/staffMember/product?id=x&vid=x
		staffMember.DELETE("/variant", handlers.DeleteProductVariant)
		// staffMember.POST("/auth-with-password", AuthWithPassword)
		// staffMember.POST("/request-password-reset", RequestPasswordReset)
		// staffMember.POST("/confirm-password-reset", ConfirmPasswordReset)
		// staffMember.POST("/auth-refresh", AuthRefresh)
		// staffMember.GET("", list) // get StaffMember List
		// staffMember.POST("", create) // create new staffMember
		// staffMember.GET("/:id", view) // view a staffMember detail
		// staffMember.PATCH("/:id", update) // update the staffMember details
		// staffMember.DELETE("/:id", delete) // delete the staffMember
	}
	return router
}

func hashStaffMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashNoStaffMember, err := auth.HashStaffMember()
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		if !hashNoStaffMember {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Contact StaffMember for signup",
			})
			ctx.Abort()
			return
		}
	}
}
