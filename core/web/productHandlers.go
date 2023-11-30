// Todo : remove the useless data

package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	database "github.com/jyotirmoydotdev/openfy/Database"
)

type RequestProduct struct {
	Handle          string   `json:"handle"`
	Description     string   `json:"description"` // Optional
	Status          bool     `json:"status"`
	Tags            []string `json:"tags"`            // Optional
	Collections     []string `json:"collections"`     // Optional
	ProductCategory string   `json:"productCategory"` // Optional
	Options         []struct {
		Name     string   `json:"name"`
		Position int      `json:"position"`
		Values   []string `json:"values"`
	} `json:"options"`
	Variants []struct {
		Price          float64 `json:"price"`
		CompareAtPrice float64 `json:"compareAtPrice"` // Optional
		CostPerItem    float64 `json:"costPerItem"`    // Optional
		Taxable        bool    `json:"taxable"`        // Optional
		Barcode        string  `json:"barcode"`        // Optional
		SKU            string  `json:"sku"`            // Optional
		Weight         struct {
			Value float64 `json:"value"`
			Uint  string  `json:"uint"`
		} `json:"weight"` // Optional
		SelectedOptions []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"selectedOptions"`
		Inventory struct {
			Available int `json:"available"`
		} `json:"inventory"`
	} `json:"variants"`
}

func Create(ctx *gin.Context) {
	var product RequestProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var productDatabase database.Product
	database.ProductMapID = make(map[string]database.Product)
	productDatabase.ID = generateProductID()
	err := copier.Copy(&productDatabase, &product)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Handle Should not be empty
	if product.Handle == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Title requied",
		})
		return
	}
	if len(product.Options) != 0 {
		TotalVariants := 1
		for i := range product.Options {
			TotalVariants *= len(product.Options[i].Values)
		}
		productDatabase.TotalVariants = TotalVariants
		productDatabase.HasOnlyDefaultVariant = false
	} else {
		productDatabase.TotalVariants = 1
		productDatabase.HasOnlyDefaultVariant = true
	}
	if len(product.Variants) != productDatabase.TotalVariants {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough variants",
		})
		return
	}
	for i := range product.Variants {
		if product.Variants[i].Price == 0.0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product does not have Price",
			})
			return
		}
		if product.Variants[i].SKU != "" {
			productDatabase.HasSKUs = true
		}
		if product.Variants[i].Barcode != "" {
			productDatabase.HasBarcodes = true
		}
		productDatabase.TotalInventory += product.Variants[i].Inventory.Available
		if product.Variants[i].CostPerItem != 0.0 {
			productDatabase.Variants[i].Profit = product.Variants[i].Price - product.Variants[i].CostPerItem
			productDatabase.Variants[i].Margin = (((product.Variants[i].Price - product.Variants[i].Price) / product.Variants[i].Price) * 100)
		}
		if product.Variants[i].Weight.Value != 0.0 {
			productDatabase.Variants[i].RequiresShipping = true
		}
		productDatabase.Variants[i].Inventory.OnHand = product.Variants[i].Inventory.Available
	}
	database.ProductList = append(database.ProductList, productDatabase)
	database.ProductMapID[productDatabase.ID] = productDatabase
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product add successfully",
	})
}
func GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, database.ProductList)
}
func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedProduct RequestProduct
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	updatedproductDatabase, ok := database.ProductMapID[id]
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}
	err := copier.Copy(&updatedproductDatabase, &updatedProduct)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Handle Should not be empty
	if updatedProduct.Handle == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Title requied",
		})
		return
	}
	if len(updatedProduct.Options) != 0 {
		TotalVariants := 1
		for i := range updatedProduct.Options {
			TotalVariants *= len(updatedProduct.Options[i].Values)
		}
		updatedproductDatabase.TotalVariants = TotalVariants
		updatedproductDatabase.HasOnlyDefaultVariant = false
	} else {
		updatedproductDatabase.TotalVariants = 1
		updatedproductDatabase.HasOnlyDefaultVariant = true
	}
	if len(updatedProduct.Variants) != updatedproductDatabase.TotalVariants {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough variants",
		})
		return
	}
	for i := range updatedProduct.Variants {
		if updatedProduct.Variants[i].Price == 0.0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product does not have Price",
			})
			return
		}
		if updatedProduct.Variants[i].SKU != "" {
			updatedproductDatabase.HasSKUs = true
		}
		if updatedProduct.Variants[i].Barcode != "" {
			updatedproductDatabase.HasBarcodes = true
		}
		updatedproductDatabase.TotalInventory += updatedProduct.Variants[i].Inventory.Available
		if updatedProduct.Variants[i].CostPerItem != 0.0 {
			updatedproductDatabase.Variants[i].Profit = updatedProduct.Variants[i].Price - updatedProduct.Variants[i].CostPerItem
			updatedproductDatabase.Variants[i].Margin = (((updatedProduct.Variants[i].Price - updatedProduct.Variants[i].Price) / updatedProduct.Variants[i].Price) * 100)
		}
		if updatedProduct.Variants[i].Weight.Value != 0.0 {
			updatedproductDatabase.Variants[i].RequiresShipping = true
		}
		updatedproductDatabase.Variants[i].Inventory.OnHand = updatedProduct.Variants[i].Inventory.Available
	}
	for i := range database.ProductList {
		if database.ProductList[i].ID == id {
			database.ProductList[i] = updatedproductDatabase
			break
		}
	}
	database.ProductMapID[id] = updatedproductDatabase
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product updated successfully",
	})
}
func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var resetProductDetails database.Product
	database.ProductMapID[id] = resetProductDetails
	for i, p := range database.ProductList {
		if p.ID == id {
			database.ProductList = append(database.ProductList[:i], database.ProductList[i+1:]...)
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
	ProductDetail, ok := database.ProductMapID[id]
	if ok {
		ctx.JSON(http.StatusOK, ProductDetail)
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Product not Found",
		})
		return
	}
}
func generateProductID() string {
	database.ProductIDCounter++
	return fmt.Sprintf("P%d", database.ProductIDCounter)
}
