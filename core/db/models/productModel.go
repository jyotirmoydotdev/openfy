package models

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type Product struct {
	ID                    uint      `gorm:"column:id;primaryKey"`
	Handle                string    `gorm:"column:handle"`
	Description           string    `gorm:"column:description"`
	Status                bool      `gorm:"column:status"`
	TotalVariants         int       `gorm:"column:totalVariants"`
	TotalInventory        int       `gorm:"column:totalInventory"`
	HasOnlyDefaultVariant bool      `gorm:"column:hasOnlyDefaultVariant"`
	OnlineStoreURL        string    `gorm:"column:onlineStoreUrl"`
	HasSKUs               bool      `gorm:"column:hasSkus"`
	HasBarcodes           bool      `gorm:"column:hasBarcodes"`
	SKURequired           bool      `gorm:"column:skuRequired"`
	Tags                  string    `gorm:"column:tags"`
	Collections           string    `gorm:"column:collections"`
	ProductCategory       string    `gorm:"column:productCategory"`
	SEOTitle              string    `gorm:"column:SEOTitle"`
	SEODescription        string    `gorm:"column:SEODescription"`
	Options               []Option  `gorm:"foreignKey:ProductID"`
	Variants              []Variant `gorm:"foreignKey:ProductID"`
}

type Option struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID int    `gorm:"column:product_id"`
	Name      string `gorm:"column:name"`
	Position  int    `gorm:"column:position"`
	Values    string `gorm:"column:values"`
}

type Variant struct {
	ID                                 uint             `gorm:"primaryKey"`
	ProductID                          int              `gorm:"column:product_id"`
	Price                              float64          `gorm:"column:price"`
	CompareAtPrice                     float64          `gorm:"column:compareAtPrice"`
	CostPerItem                        float64          `gorm:"column:costPerItem"`
	Taxable                            bool             `gorm:"column:taxable"`
	Profit                             float64          `gorm:"column:profit"`
	Margin                             float64          `gorm:"column:margin"`
	Barcode                            string           `gorm:"column:barcode"`
	SKU                                string           `gorm:"column:sku"`
	RequiresShipping                   bool             `gorm:"column:requiresShipping"`
	WeightValue                        float64          `gorm:"column:weightValue"`
	WeightUnit                         string           `gorm:"column:weightUnit"`
	SelectedOptions                    []SelectedOption `gorm:"foreignKey:VariantID"`
	InventoryAvailable                 int              `gorm:"column:inventoryAvailable"`
	InventoryCommitted                 int              `gorm:"column:inventoryCommitted"`
	InventoryOnHand                    int              `gorm:"column:inventoryOnHand"`
	InventoryUnavailableDamaged        int              `gorm:"column:inventoryUnavailableDamaged"`
	InventoryUnavailableQualityControl int              `gorm:"column:inventoryUnavailableQualityControl"`
	InventoryUnavailableSafetyStock    int              `gorm:"column:inventoryUnavailableSafetyStock"`
	InventoryUnavailableOther          int              `gorm:"column:inventoryUnavailableOther"`
}

type SelectedOption struct {
	ID        uint   `gorm:"primaryKey"`
	VariantID int    `gorm:"column:variant_id"`
	Name      string `gorm:"column:name"`
	Value     string `gorm:"column:value"`
}

var ProductMapID map[string]Product
var ProductList []Product
var ProductIDCounter int

// var Tags []string
// var Collections []string

type ProductModel struct {
	db *gorm.DB
}

func NewProductModel(db *gorm.DB) *ProductModel {
	return &ProductModel{db: db}
}
func (pd *ProductModel) Save(product *Product) error {
	return pd.db.Create(&product).Error
}

func (pd *ProductModel) GetProduct(id string) (*Product, error) {
	var existingProduct Product
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").First(&existingProduct, uintID).Error; err != nil {
		return nil, err
	}
	return &existingProduct, nil
}

// TODO : fix update product
func (pd *ProductModel) Update(id string, product *Product) error {
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errors.New("invalid ID format")
	}
	// var existingProduct Product
	// if err := pd.db.First(&existingProduct, "id = ?", uintID).Error; err != nil {
	// 	return errors.New("product not found")
	// }
	product.ID = uint(uintID)
	return pd.db.Model(&Product{}).Where("id = ?", uint(uintID)).Updates(product).Error
}

func (pd *ProductModel) Delete(id string, product *Product) error {
	return nil
}
