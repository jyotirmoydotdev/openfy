package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
)

type RequestProduct struct {
	Handle          string     `json:"handle"`
	Description     string     `json:"description"`
	Status          bool       `json:"status"`
	Tags            []string   `json:"product_tags"`
	Collections     []string   `json:"product_collections"`
	ProductCategory string     `json:"productCategory"`
	Options         []Options  `json:"options"`
	Variants        []Variants `json:"variants"`
}
type Options struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Values   []string `json:"values"`
}
type Variants struct {
	Price           float64 `json:"price"`
	CompareAtPrice  float64 `json:"compareAtPrice"`
	CostPerItem     float64 `json:"costPerItem"`
	Taxable         bool    `json:"taxable"`
	Barcode         string  `json:"barcode"`
	SKU             string  `json:"sku"`
	WeightValue     float64 `json:"weightValue"`
	WeightUnit      string  `json:"weightUnit"`
	SelectedOptions []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"selectedOptions"`
	InventoryAvailable int `json:"inventoryAvailable"`
}

func NewRequestProductHandlers() *RequestProduct {
	return &RequestProduct{}
}
func concatenateStrings(slice []string) string {
	return strings.Join(slice, ",")
}

// Create a new product
// Expected : 200
func (rp *RequestProduct) Create(ctx *gin.Context) {
	var product RequestProduct

	// Bind the request body to the RequestProduct struct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var productDatabase models.Product

	// Copy the values from the request to the database model
	err := copier.Copy(&productDatabase, &product)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the slice of strings to a single string
	productDatabase.Tags = concatenateStrings(product.Tags)
	productDatabase.Collections = concatenateStrings(product.Collections)
	for i, v := range product.Options {
		productDatabase.Options[i].Values = concatenateStrings(v.Values)
	}

	// Handle Should not be empty
	if productDatabase.Handle == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Title requied",
		})
		return
	}

	// Status shoule be True by default
	productDatabase.Status = true

	// Format the data related to variants
	for i := range productDatabase.Variants {
		if productDatabase.Variants[i].Price == 0.0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product does not have Price",
			})
			return
		}
		if productDatabase.Variants[i].SKU != "" {
			productDatabase.HasSKUs = true
		}
		if productDatabase.Variants[i].Barcode != "" {
			productDatabase.HasBarcodes = true
		}
		productDatabase.TotalInventory += productDatabase.Variants[i].InventoryAvailable
		if productDatabase.Variants[i].CostPerItem != 0.0 {
			productDatabase.Variants[i].Profit = productDatabase.Variants[i].Price - productDatabase.Variants[i].CostPerItem
			productDatabase.Variants[i].Margin = (((productDatabase.Variants[i].Price - productDatabase.Variants[i].Price) / productDatabase.Variants[i].Price) * 100)
		}
		if productDatabase.Variants[i].WeightValue != 0.0 {
			productDatabase.Variants[i].RequiresShipping = true
		}
		productDatabase.Variants[i].InventoryOnHand = productDatabase.Variants[i].InventoryAvailable
	}

	// Connect to the database
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)

	// Save the product to the database
	if err := productModel.Save(&productDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Return a success message
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product add successfully",
	})
}

//	func (rp *RequestProduct) Update(ctx *gin.Context) {
//		id := ctx.Param("id")
//		var updatedProduct RequestProduct
//		if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"error": err.Error(),
//			})
//		}
//		updatedproductDatabase, ok := models.ProductMapID[id]
//		if !ok {
//			ctx.JSON(http.StatusNotFound, gin.H{
//				"error": "Product not found",
//			})
//			return
//		}
//		err := copier.Copy(&updatedproductDatabase, &updatedProduct)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return
//		}
//		// Handle Should not be empty
//		if updatedProduct.Handle == "" {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"error": "Title requied",
//			})
//			return
//		}
//		if len(updatedProduct.Options) != 0 {
//			TotalVariants := 1
//			for i := range updatedProduct.Options {
//				TotalVariants *= len(updatedProduct.Options[i].Values)
//			}
//			updatedproductDatabase.TotalVariants = TotalVariants
//			updatedproductDatabase.HasOnlyDefaultVariant = false
//		} else {
//			updatedproductDatabase.TotalVariants = 1
//			updatedproductDatabase.HasOnlyDefaultVariant = true
//		}
//		if len(updatedProduct.Variants) != updatedproductDatabase.TotalVariants {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"error": "Not enough variants",
//			})
//			return
//		}
//		for i := range updatedProduct.Variants {
//			if updatedProduct.Variants[i].Price == 0.0 {
//				ctx.JSON(http.StatusBadRequest, gin.H{
//					"error": "Product does not have Price",
//				})
//				return
//			}
//			if updatedProduct.Variants[i].SKU != "" {
//				updatedproductDatabase.HasSKUs = true
//			}
//			if updatedProduct.Variants[i].Barcode != "" {
//				updatedproductDatabase.HasBarcodes = true
//			}
//			updatedproductDatabase.TotalInventory += updatedProduct.Variants[i].InventoryAvailable
//			if updatedProduct.Variants[i].CostPerItem != 0.0 {
//				updatedproductDatabase.Variants[i].Profit = updatedProduct.Variants[i].Price - updatedProduct.Variants[i].CostPerItem
//				updatedproductDatabase.Variants[i].Margin = (((updatedProduct.Variants[i].Price - updatedProduct.Variants[i].Price) / updatedProduct.Variants[i].Price) * 100)
//			}
//			if updatedProduct.Variants[i].WeightValue != 0.0 {
//				updatedproductDatabase.Variants[i].RequiresShipping = true
//			}
//			updatedproductDatabase.Variants[i].InventoryOnHand = updatedProduct.Variants[i].InventoryAvailable
//		}
//		for i := range models.ProductList {
//			if string(models.ProductList[i].ID) == id {
//				models.ProductList[i] = updatedproductDatabase
//				break
//			}
//		}
//		models.ProductMapID[id] = updatedproductDatabase
//		ctx.JSON(http.StatusOK, gin.H{
//			"status": "Product updated successfully",
//		})
//	}
func (rp *RequestProduct) GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.ProductList)
}
func (rp *RequestProduct) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var resetProductDetails models.Product
	models.ProductMapID[id] = resetProductDetails
	for i, p := range models.ProductList {
		if string(p.ID) == id {
			models.ProductList = append(models.ProductList[:i], models.ProductList[i+1:]...)
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
func (rp *RequestProduct) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	ProductDetail, ok := models.ProductMapID[id]
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
	models.ProductIDCounter++
	return fmt.Sprintf("P%d", models.ProductIDCounter)
}
