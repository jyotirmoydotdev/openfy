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
	ID              uint       `json:"id"`
	Handle          string     `json:"handle"`
	Description     string     `json:"description"`
	Status          bool       `json:"status"`
	Tags            []string   `json:"tags"`
	Collections     []string   `json:"collections"`
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

// Create a new product
// Expected : 200
func (rp *RequestProduct) Create(ctx *gin.Context) {
	var createProduct RequestProduct

	// Bind the request body to the RequestProduct struct
	if err := ctx.ShouldBindJSON(&createProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productData, err := validateProduct(createProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
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
	if err := productModel.Save(productData); err != nil {
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
func (rp *RequestProduct) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id is missing",
		})
	}
	var updateProduct RequestProduct

	// Bind the request body to the RequestProduct struct
	if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: replace the _ with productData and save it to database
	productDatabase, err := validateProduct(updateProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
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
	if err := productModel.Update(id, productDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

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
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)
	product, err := productModel.GetProduct(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func concatenateStrings(slice []string) string {
	return strings.Join(slice, ",")
}
func validateProduct(product RequestProduct) (*models.Product, error) {
	var productDatabase models.Product
	// Copy the values from the request to the database model
	err := copier.Copy(&productDatabase, &product)
	if err != nil {
		return nil, err
	}

	// Convert the slice of strings to a single string
	productDatabase.Tags = concatenateStrings(product.Tags)
	productDatabase.Collections = concatenateStrings(product.Collections)
	for i, v := range product.Options {
		productDatabase.Options[i].Values = concatenateStrings(v.Values)
	}

	// Handle Should not be empty
	if productDatabase.Handle == "" {
		return nil, fmt.Errorf("title required")
	}

	// Status shoule be True by default
	productDatabase.Status = true

	// Add Total Variants
	productDatabase.TotalVariants = 1
	for i := range product.Options {
		productDatabase.TotalVariants *= len(product.Options[i].Values)
	}
	if productDatabase.TotalVariants != len(productDatabase.Variants) {
		return nil, fmt.Errorf("not enough variants")
	}
	if productDatabase.TotalVariants == 1 {
		productDatabase.HasOnlyDefaultVariant = true
	} else {
		productDatabase.HasOnlyDefaultVariant = false
	}

	// Format the data related to variants
	for i := range productDatabase.Variants {
		if productDatabase.Variants[i].Price == 0.0 {
			return nil, fmt.Errorf("product does not have Price")
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
			productDatabase.Variants[i].Margin = (((productDatabase.Variants[i].Price - productDatabase.Variants[i].CostPerItem) / productDatabase.Variants[i].Price) * 100)
		}
		if productDatabase.Variants[i].WeightValue != 0.0 {
			productDatabase.Variants[i].RequiresShipping = true
		}
		productDatabase.Variants[i].InventoryOnHand = productDatabase.Variants[i].InventoryAvailable
	}
	return &productDatabase, nil
}
