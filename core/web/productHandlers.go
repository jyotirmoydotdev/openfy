// Todo : remove the useless data

package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Product struct {
	ID                    string   `json:"id"`
	Handle                string   `json:"handle"`
	Description           string   `json:"description"`
	Status                bool     `json:"status"`
	TotalVariants         int      `json:"totalVariants"`
	TotalInventory        int      `json:"totalInventory"`
	HasOnlyDefaultVariant bool     `json:"hasOnlyDefaultVariant"`
	OnlineStoreURL        string   `json:"onlineStoreUrl"`
	HasSKUs               bool     `json:"hasSkus"`
	HasBarcodes           bool     `json:"hasBarcodes"`
	SKURequired           bool     `json:"skuRequired"`
	Tags                  []string `json:"tags"`
	Collections           []string `json:"collections"`
	ProductCategory       string   `json:"productCategory"`
	Options               []struct {
		ID       string   `json:"id"`
		Name     string   `json:"name"`
		Position int      `json:"position"`
		Values   []string `json:"values"`
	} `json:"options"`
	SEO struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"seo"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Price            float64 `json:"price"`
	CompareAtPrice   float64 `json:"compareAtPrice"`
	CostPerItem      float64 `json:"costPerItem"`
	Taxable          bool    `json:"taxable"`
	Profit           float64 `json:"profit"`
	Margin           float64 `json:"margin"`
	Barcode          string  `json:"barcode"`
	SKU              string  `json:"sku"`
	RequiresShipping bool    `json:"requiresShipping"`
	Weight           struct {
		Value float64 `json:"value"`
		Uint  string  `json:"uint"`
	} `json:"weight"`
	SelectedOptions []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"selectedOptions"`
	Inventory struct {
		Available   int `json:"available"`
		Committed   int `json:"committed"`
		OnHand      int `json:"onHand"`
		Unavailable struct {
			Damaged struct {
				Quantity int `json:"quantity"`
			} `json:"damaged"`
			QualityControl struct {
				Quantity int `json:"quantity"`
			} `json:"qualityControl"`
			SafetyStock struct {
				Quantity int `json:"quantity"`
			} `json:"safetyStock"`
			Other struct {
				Quantity int `json:"quantity"`
			} `json:"other"`
		} `json:"unavailable"`
	} `json:"inventory"`
}

var ProductMapID map[string]Product
var ProductList []Product
var productIDCounter int

func Create(ctx *gin.Context) {
	var product struct {
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
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var Database Product
	ProductMapID = make(map[string]Product)
	Database.ID = generateProductID()
	err := copier.Copy(&Database, &product)
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
	// Check if there is vairants and inventory
	// If lenght of product.Variants != Database.TotalVariants
	// Return error
	if len(product.Options) != 0 {
		TotalVariants := 1
		for i := range product.Options {
			TotalVariants *= len(product.Options[i].Values)
		}
		Database.TotalVariants = TotalVariants
		Database.HasOnlyDefaultVariant = false
	} else {
		Database.TotalVariants = 1
		Database.HasOnlyDefaultVariant = true
	}
	if len(product.Variants) != Database.TotalVariants {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough variants",
		})
		return
	}
	// save Total Inventory to database
	// check if it has sku and barcode
	for i := range product.Variants {
		if product.Variants[i].Price == 0.0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product does not have Price",
			})
			return
		}
		if product.Variants[i].SKU != "" {
			Database.HasSKUs = true
		}
		if product.Variants[i].Barcode != "" {
			Database.HasBarcodes = true
		}
		Database.TotalInventory += product.Variants[i].Inventory.Available
		if product.Variants[i].CostPerItem != 0.0 {
			Database.Variants[i].Profit = product.Variants[i].Price - product.Variants[i].CostPerItem
			Database.Variants[i].Margin = (((product.Variants[i].Price - product.Variants[i].Price) / product.Variants[i].Price) * 100)
		}
		if product.Variants[i].Weight.Value != 0.0 {
			Database.Variants[i].RequiresShipping = true
		}
		Database.Variants[i].Inventory.OnHand = product.Variants[i].Inventory.Available
	}
	ProductList = append(ProductList, Database)
	ProductMapID[Database.ID] = Database
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
	ProductDetail, ok := ProductMapID[id]
	if ok {
		ctx.JSON(http.StatusOK, ProductDetail)
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Product not Found",
		})
	}
}

func generateProductID() string {
	productIDCounter++
	return fmt.Sprintf("P%d", productIDCounter)
}
