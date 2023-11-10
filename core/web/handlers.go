package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Media            []string `json:"media"`
	Price            int      `json:"price"`
	Compare_At_Price int      `json:"compareatprice"`
	Tax              bool     `json:"tax"`
	Cost_Per_Item    int      `json:"costperitem"`
}

var ProductList []Product

func Create(ctx *gin.Context) {
	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if product.Title == "" || product.Price <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Title and Price are requied",
		})
		return
	}
	ProductList = append(ProductList, product)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product add successfully",
	})
}

func GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ProductList)
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateProduct Product
	if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	for i, p := range ProductList {
		if p.ID == id {
			ProductList[i] = updateProduct
			ctx.JSON(http.StatusOK, gin.H{
				"Status": "Product updated successfully",
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Product not found",
	})
}
func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	for i, p := range ProductList {
		if p.ID == id {
			ProductList = append(ProductList[:i], ProductList[i+1:]...)
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Product deleted successfully",
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Product not found",
	})
}

func GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	for _, p := range ProductList {
		if p.ID == id {
			ctx.JSON(http.StatusOK, p)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Product not Found",
	})
}
