package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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

// Create a new product
// Expected : 200
func Create(ctx *gin.Context) {
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
func Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
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
func GetAllProducts(ctx *gin.Context) {
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)
	allProduct, err := productModel.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": allProduct,
	})
}
func DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)
	err = productModel.DeleteProduct(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "product deleted succesfully",
	})
}
func DeleteProductVariant(ctx *gin.Context) {
	id, err1 := strconv.Atoi(ctx.Query("id"))
	vid, err2 := strconv.Atoi(ctx.Query("vid"))
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)
	err = productModel.DeleteProductVariant(id, vid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "varient deleted succesfully",
	})
}
func GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
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

func GetAllActiveProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	productdbInstance, err := db.GetProductDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	productModel := models.NewProductModel(productdbInstance)
	products, err := productModel.GetPaginatedActiveProducts(page, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var userProduct []UserProduct
	for i := range products {
		userProduct = append(userProduct, ConvertToUserProduct(products[i]))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": userProduct,
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

type UserOption struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Values string `json:"values"`
}

type UserVariant struct {
	ID                 uint    `json:"id"`
	Price              float64 `json:"price"`
	CompareAtPrice     float64 `json:"compareAtPrice"`
	CostPerItem        float64 `json:"costPerItem"`
	Taxable            bool    `json:"taxable"`
	RequiresShipping   bool    `json:"requiresShipping"`
	WeightValue        float64 `json:"weightValue"`
	WeightUnit         string  `json:"weightUnit"`
	InventoryAvailable int     `json:"inventoryAvailable"`
}

type UserProduct struct {
	ID              uint          `json:"id"`
	Handle          string        `json:"handle"`
	Description     string        `json:"description"`
	TotalVariants   int           `json:"totalVariants"`
	TotalInventory  int           `json:"totalInventory"`
	OnlineStoreURL  string        `json:"onlineStoreUrl"`
	HasSKUs         bool          `json:"hasSkus"`
	HasBarcodes     bool          `json:"hasBarcodes"`
	ProductCategory string        `json:"productCategory"`
	Options         []UserOption  `json:"options"`
	Variants        []UserVariant `json:"variants"`
}

func ConvertToUserProduct(adminProduct models.Product) UserProduct {
	userProduct := UserProduct{
		ID:              adminProduct.ID,
		Handle:          adminProduct.Handle,
		Description:     adminProduct.Description,
		TotalVariants:   adminProduct.TotalVariants,
		TotalInventory:  adminProduct.TotalInventory,
		OnlineStoreURL:  adminProduct.OnlineStoreURL,
		HasSKUs:         adminProduct.HasSKUs,
		HasBarcodes:     adminProduct.HasBarcodes,
		ProductCategory: adminProduct.ProductCategory,
		Options:         make([]UserOption, len(adminProduct.Options)),
		Variants:        make([]UserVariant, len(adminProduct.Variants)),
	}

	// Convert Options
	for i, adminOption := range adminProduct.Options {
		userProduct.Options[i] = UserOption{
			ID:     adminOption.ID,
			Name:   adminOption.Name,
			Values: adminOption.Values,
		}
	}

	// Convert Variants
	for i, adminVariant := range adminProduct.Variants {
		userProduct.Variants[i] = UserVariant{
			ID:                 adminVariant.ID,
			Price:              adminVariant.Price,
			CompareAtPrice:     adminVariant.CompareAtPrice,
			CostPerItem:        adminVariant.CostPerItem,
			Taxable:            adminVariant.Taxable,
			RequiresShipping:   adminVariant.RequiresShipping,
			WeightValue:        adminVariant.WeightValue,
			WeightUnit:         adminVariant.WeightUnit,
			InventoryAvailable: adminVariant.InventoryAvailable,
		}
	}

	return userProduct
}
