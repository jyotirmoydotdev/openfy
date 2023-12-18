package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
)

type RequestProduct struct {
	Handle          string     `json:"handle"`
	Description     string     `json:"description"` // Optional
	Status          bool       `json:"status"`
	Tags            []string   `json:"tags"`            // Optional
	Collections     []string   `json:"collections"`     // Optional
	ProductCategory string     `json:"productCategory"` // Optional
	Options         []Options  `json:"options"`
	Variants        []Variants `json:"variants"`
}
type Options struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Values   []string `json:"values"`
}
type Variants struct {
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
}

func NewRequestProductHandlers() *RequestProduct {
	return &RequestProduct{}
}

func (rp *RequestProduct) Create(ctx *gin.Context) {
	var product RequestProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// models.ProductMapID = make(map[string]models.Product)

	// ------------------------------------------------------------
	var productDatabase models.Product
	productDatabase.ID = generateProductID() // <-----------------
	err := copier.Copy(&productDatabase, &product)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// ------------------------------------------------------------

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

	// CHECK -> models.ProductList = append(models.ProductList, productDatabase)
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	}
	productModel := models.NewProductModel(productdbInstance)
	if err := productModel.Save(&productDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// models.ProductMapID[productDatabase.ID] = productDatabase
	// for _, v := range productDatabase.Tags {
	// 	if slices.Contains(models.Tags, v) {
	// 		continue
	// 	} else {
	// 		models.Tags = append(models.Tags, v)
	// 	}
	// }
	// for _, v := range productDatabase.Collections {
	// 	if slices.Contains(models.Collections, v) {
	// 		continue
	// 	} else {
	// 		models.Collections = append(models.Collections, v)
	// 	}
	// }
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product add successfully",
	})
}

func (rp *RequestProduct) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedProduct RequestProduct
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	updatedproductDatabase, ok := models.ProductMapID[id]
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
	for i := range models.ProductList {
		if models.ProductList[i].ID == id {
			models.ProductList[i] = updatedproductDatabase
			break
		}
	}
	models.ProductMapID[id] = updatedproductDatabase
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Product updated successfully",
	})
}
func (rp *RequestProduct) GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.ProductList)
}
func (rp *RequestProduct) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var resetProductDetails models.Product
	models.ProductMapID[id] = resetProductDetails
	for i, p := range models.ProductList {
		if p.ID == id {
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
