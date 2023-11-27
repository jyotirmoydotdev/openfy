// Todo : remove the useless data

package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID                    string `json:"id"`
	Handle                string `json:"handle"`
	Description           string `json:"description"`
	Status                string `json:"status"`
	TotalVariants         int    `json:"totalVariants"`
	TotalInventory        int    `json:"totalInventory"`
	HasOnlyDefaultVariant bool   `json:"hasOnlyDefaultVariant"`
	OnlineStoreURL        string `json:"onlineStoreUrl"`
	HasSKUs               bool   `json:"hasSkus"`
	HasBarcodes           bool   `json:"hasBarcodes"`
	SKURequired           bool   `json:"skuRequired"`
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
	ProductCategory string    `json:"productCategory"`
	Variants        []Variant `json:"variants"`
}

type Variant struct {
	Price            float64 `json:"price"`
	CompareAtPrice   float64 `json:"compareAtPrice"`
	Taxable          bool    `json:"taxable"`
	Profit           float64 `json:"profit"`
	CostPerItem      float64 `json:"costPerItem"`
	Margin           float64 `json:"margin"`
	RequiresShipping bool    `json:"requiresShipping"`
	Weight           float64 `json:"weight"`
	Barcode          string  `json:"barcode"`
	SKU              string  `json:"sku"`
	SelectedOptions  []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"selectedOptions"`
	Inventory struct {
		Available             int    `json:"available"`
		Committed             int    `json:"committed"`
		OnHand                int    `json:"onHand"`
		CanDeactivate         bool   `json:"canDeactivate"`
		DeactivationAlertHtml string `json:"deactivationAlertHtml"`
		Incoming              int    `json:"incoming"`
		Unavailable           []struct {
			Quantity int    `json:"quantity"`
			Name     string `json:"name"`
		} `json:"unavailable"`
	} `json:"inventory"`
}

var ProductList []Product

var productIDCounter int

func Create(ctx *gin.Context) {
	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.ID = generateProductID()

	if product.Handle == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Title requied",
		})
		return
	}
	if product.HasOnlyDefaultVariant {
		if product.Variants[0].Price <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Price requied",
			})
			return
		}
	} else {
		for i := range product.Variants {
			if product.Variants[i].Price <= 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Price requied",
				})
				return
			}
		}
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

func generateProductID() string {
	productIDCounter++
	return fmt.Sprintf("P%d", productIDCounter)
}
