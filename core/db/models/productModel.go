package models

import "gorm.io/gorm"

type Product struct {
	ID                    string    `gorm:"column:id;primaryKey"`
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
	Tags                  []string  `gorm:"column:product_tags;type:json"`
	Collections           []string  `gorm:"column:product_collections;type:json"`
	ProductCategory       string    `gorm:"column:productCategory"`
	SEO                   SEO       `gorm:"embedded;embeddedPrefix:seo_"`
	Options               []Option  `gorm:"foreignKey:ProductID"`
	Variants              []Variant `gorm:"foreignKey:ProductID"`
}

type Option struct {
	ID        string   `gorm:"column:id;primaryKey"`
	ProductID string   `gorm:"column:product_id"`
	Name      string   `gorm:"column:name"`
	Position  int      `gorm:"column:position"`
	Values    []string `gorm:"many2many:values;type:json"`
}

type SEO struct {
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

type Variant struct {
	ID               string           `gorm:"column:id"`
	ProductID        string           `gorm:"column:product_id"`
	Price            float64          `gorm:"column:price"`
	CompareAtPrice   float64          `gorm:"column:compareAtPrice"`
	CostPerItem      float64          `gorm:"column:costPerItem"`
	Taxable          bool             `gorm:"column:taxable"`
	Profit           float64          `gorm:"column:profit"`
	Margin           float64          `gorm:"column:margin"`
	Barcode          string           `gorm:"column:barcode"`
	SKU              string           `gorm:"column:sku"`
	RequiresShipping bool             `gorm:"column:requiresShipping"`
	Weight           Weight           `gorm:"embedded;embeddedPrefix:weight_"`
	SelectedOptions  []SelectedOption `gorm:"foreignKey:VariantID"`
	Inventory        Inventory        `gorm:"embedded;embeddedPrefix:inventory_"`
}

type Weight struct {
	Value float64 `gorm:"column:value"`
	Unit  string  `gorm:"column:uint"`
}

type SelectedOption struct {
	VariantID string `gorm:"column:variant_id"`
	Name      string `gorm:"column:name"`
	Value     string `gorm:"column:value"`
}

type Inventory struct {
	Available   int         `gorm:"column:available"`
	Committed   int         `gorm:"column:committed"`
	OnHand      int         `gorm:"column:onHand"`
	Unavailable Unavailable `gorm:"embedded;embeddedPrefix:unavailable_"`
}

type Unavailable struct {
	Damaged        Damaged        `gorm:"embedded;embeddedPrefix:damaged_"`
	QualityControl QualityControl `gorm:"embedded;embeddedPrefix:qc_"`
	SafetyStock    SafetyStock    `gorm:"embedded;embeddedPrefix:safety_"`
	Other          Other          `gorm:"embedded;embeddedPrefix:other_"`
}

type Damaged struct {
	Quantity int `gorm:"column:quantity"`
}

type QualityControl struct {
	Quantity int `gorm:"column:quantity"`
}

type SafetyStock struct {
	Quantity int `gorm:"column:quantity"`
}

type Other struct {
	Quantity int `gorm:"column:quantity"`
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
	return pd.db.Create(product).Error
}
