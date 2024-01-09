package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Product struct {
	ID                    uint      `gorm:"column:id;primaryKey" json:"id"`
	Handle                string    `gorm:"column:handle" json:"handle"`
	Description           string    `gorm:"column:description" json:"description"`
	Status                bool      `gorm:"column:status" json:"status"`
	TotalVariants         int       `gorm:"column:totalVariants" json:"totalVariants"`
	TotalInventory        int       `gorm:"column:totalInventory" json:"totalInventory"`
	HasOnlyDefaultVariant bool      `gorm:"column:hasOnlyDefaultVariant" json:"hasOnlyDefaultVariant"`
	OnlineStoreURL        string    `gorm:"column:onlineStoreUrl" json:"onlineStoreUrl"`
	HasSKUs               bool      `gorm:"column:hasSkus" json:"hasSkus"`
	HasBarcodes           bool      `gorm:"column:hasBarcodes" json:"hasBarcodes"`
	SKURequired           bool      `gorm:"column:skuRequired" json:"skuRequired"`
	Tags                  string    `gorm:"column:tags" json:"tags"`
	Collections           string    `gorm:"column:collections" json:"collections"`
	ProductCategory       string    `gorm:"column:productCategory" json:"productCategory"`
	SEOTitle              string    `gorm:"column:SEOTitle" json:"SEOTitle"`
	SEODescription        string    `gorm:"column:SEODescription" json:"SEODescription"`
	Options               []Option  `gorm:"foreignKey:ProductID"`
	Variants              []Variant `gorm:"foreignKey:ProductID"`
}

type Option struct {
	ID        uint   `gorm:"primaryKey" json:"id" `
	ProductID uint   `gorm:"column:product_id" json:"product_id"  `
	Name      string `gorm:"column:name" json:"name"  `
	Position  int    `gorm:"column:position" json:"position"  `
	Values    string `gorm:"column:values" json:"values"  `
}

type Variant struct {
	ID                                 uint             `gorm:"primaryKey" json:"id"`
	ProductID                          uint             `gorm:"column:product_id" json:"product_id"`
	Price                              float64          `gorm:"column:price" json:"price"`
	CompareAtPrice                     float64          `gorm:"column:compareAtPrice" json:"compareAtPrice"`
	CostPerItem                        float64          `gorm:"column:costPerItem" json:"costPerItem"`
	Taxable                            bool             `gorm:"column:taxable" json:"taxable"`
	Profit                             float64          `gorm:"column:profit" json:"profit"`
	Margin                             float64          `gorm:"column:margin" json:"margin"`
	Barcode                            string           `gorm:"column:barcode" json:"barcode"`
	SKU                                string           `gorm:"column:sku" json:"sku"`
	RequiresShipping                   bool             `gorm:"column:requiresShipping" json:"requiresShipping"`
	WeightValue                        float64          `gorm:"column:weightValue" json:"weightValue"`
	WeightUnit                         string           `gorm:"column:weightUnit" json:"weightUnit"`
	InventoryAvailable                 int              `gorm:"column:inventoryAvailable" json:"inventoryAvailable"`
	InventoryCommitted                 int              `gorm:"column:inventoryCommitted" json:"inventoryCommitted"`
	InventoryOnHand                    int              `gorm:"column:inventoryOnHand" json:"inventoryOnHand"`
	InventoryUnavailableDamaged        int              `gorm:"column:inventoryUnavailableDamaged" json:"inventoryUnavailableDamaged"`
	InventoryUnavailableQualityControl int              `gorm:"column:inventoryUnavailableQualityControl" json:"inventoryUnavailableQualityControl"`
	InventoryUnavailableSafetyStock    int              `gorm:"column:inventoryUnavailableSafetyStock" json:"inventoryUnavailableSafetyStock"`
	InventoryUnavailableOther          int              `gorm:"column:inventoryUnavailableOther" json:"inventoryUnavailableOther"`
	SelectedOptions                    []SelectedOption `gorm:"foreignKey:VariantID"`
}

type SelectedOption struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	VariantID uint   `gorm:"column:variant_id" json:"variant_id"`
	Name      string `gorm:"column:name" json:"name"`
	Value     string `gorm:"column:value" json:"value"`
}

type ProductModel struct {
	db *gorm.DB
}

func NewProductModel(db *gorm.DB) *ProductModel {
	return &ProductModel{db: db}
}

func (pd *ProductModel) Save(product *Product) error {
	return pd.db.Create(&product).Error
}

func (pd *ProductModel) GetProduct(id int) (*Product, error) {
	var existingProduct Product
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").First(&existingProduct, uint(id)).Error; err != nil {
		return nil, err
	}
	return &existingProduct, nil
}
func (pd *ProductModel) GetAllProducts() ([]Product, error) {
	var existingProducts []Product
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").Find(&existingProducts).Error; err != nil {
		return nil, err
	}
	return existingProducts, nil
}

func (pd *ProductModel) Update(id int, updatedProduct *Product) error {
	var existingProduct Product
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").First(&existingProduct, "id = ?", uint(id)).Error; err != nil {
		return errors.New("product not found")
	}

	existingProduct.Handle = updatedProduct.Handle
	existingProduct.Description = updatedProduct.Description
	existingProduct.Status = updatedProduct.Status
	existingProduct.TotalVariants = updatedProduct.TotalVariants
	existingProduct.TotalInventory = updatedProduct.TotalInventory
	existingProduct.HasOnlyDefaultVariant = updatedProduct.HasOnlyDefaultVariant
	existingProduct.OnlineStoreURL = updatedProduct.OnlineStoreURL
	existingProduct.HasSKUs = updatedProduct.HasSKUs
	existingProduct.HasBarcodes = updatedProduct.HasBarcodes
	existingProduct.SKURequired = updatedProduct.SKURequired
	existingProduct.Tags = updatedProduct.Tags
	existingProduct.Collections = updatedProduct.Collections
	existingProduct.ProductCategory = updatedProduct.ProductCategory
	existingProduct.SEOTitle = updatedProduct.SEOTitle
	existingProduct.SEODescription = updatedProduct.SEODescription

	for _, Variant := range existingProduct.Variants {
		if err := pd.deleteSelectedOptionsForVariant(Variant.ID); err != nil {
			return err
		}
	}

	if err := pd.db.Delete(&existingProduct.Variants).Error; err != nil {
		return err
	}

	if len(existingProduct.Options) > 0 {
		if err := pd.db.Delete(&existingProduct.Options).Error; err != nil {
			return err
		}
	}

	for i, updatedOption := range updatedProduct.Options {
		if i < len(existingProduct.Options) {
			// Update existing option
			existingProduct.Options[i].Name = updatedOption.Name
			existingProduct.Options[i].Position = updatedOption.Position
			existingProduct.Options[i].Values = updatedOption.Values
		} else {
			// Create new option
			newOption := Option{
				ProductID: existingProduct.ID,
				Name:      updatedOption.Name,
				Position:  updatedOption.Position,
				Values:    updatedOption.Values,
			}
			existingProduct.Options = append(existingProduct.Options, newOption)
		}
	}

	// Update or create Variants
	for i, updatedVariant := range updatedProduct.Variants {
		if i < len(existingProduct.Variants) {
			// Update existing variant
			existingProduct.Variants[i].Price = updatedVariant.Price
			existingProduct.Variants[i].CompareAtPrice = updatedVariant.CompareAtPrice
			existingProduct.Variants[i].CostPerItem = updatedVariant.CostPerItem
			existingProduct.Variants[i].Taxable = updatedVariant.Taxable
			existingProduct.Variants[i].Profit = updatedVariant.Profit
			existingProduct.Variants[i].Margin = updatedVariant.Margin
			existingProduct.Variants[i].Barcode = updatedVariant.Barcode
			existingProduct.Variants[i].SKU = updatedVariant.SKU
			existingProduct.Variants[i].RequiresShipping = updatedVariant.RequiresShipping
			existingProduct.Variants[i].WeightValue = updatedVariant.WeightValue
			existingProduct.Variants[i].WeightUnit = updatedVariant.WeightUnit
			existingProduct.Variants[i].InventoryAvailable = updatedVariant.InventoryAvailable
			existingProduct.Variants[i].InventoryCommitted = updatedVariant.InventoryCommitted
			existingProduct.Variants[i].InventoryOnHand = updatedVariant.InventoryOnHand
			existingProduct.Variants[i].InventoryUnavailableDamaged = updatedVariant.InventoryUnavailableDamaged
			existingProduct.Variants[i].InventoryUnavailableQualityControl = updatedVariant.InventoryUnavailableQualityControl
			existingProduct.Variants[i].InventoryUnavailableSafetyStock = updatedVariant.InventoryUnavailableSafetyStock
			existingProduct.Variants[i].InventoryUnavailableOther = updatedVariant.InventoryUnavailableOther

			// Update or create SelectedOptions
			for j, updatedSelectedOption := range updatedVariant.SelectedOptions {
				if j < len(existingProduct.Variants[i].SelectedOptions) {
					// Update existing SelectedOption
					existingProduct.Variants[i].SelectedOptions[j].VariantID = updatedSelectedOption.VariantID
					existingProduct.Variants[i].SelectedOptions[j].Name = updatedSelectedOption.Name
					existingProduct.Variants[i].SelectedOptions[j].Value = updatedSelectedOption.Value
				} else {
					// Create new SelectedOption
					newSelectedOption := SelectedOption{
						VariantID: existingProduct.Variants[i].ID,
						Name:      updatedSelectedOption.Name,
						Value:     updatedSelectedOption.Value,
					}
					existingProduct.Variants[i].SelectedOptions = append(existingProduct.Variants[i].SelectedOptions, newSelectedOption)
				}
			}
		} else {
			// Create new variant
			newVariant := Variant{
				ProductID:                          existingProduct.ID,
				Price:                              updatedVariant.Price,
				CompareAtPrice:                     updatedVariant.CompareAtPrice,
				CostPerItem:                        updatedVariant.CostPerItem,
				Taxable:                            updatedVariant.Taxable,
				Profit:                             updatedVariant.Profit,
				Margin:                             updatedVariant.Margin,
				Barcode:                            updatedVariant.Barcode,
				SKU:                                updatedVariant.SKU,
				RequiresShipping:                   updatedVariant.RequiresShipping,
				WeightValue:                        updatedVariant.WeightValue,
				WeightUnit:                         updatedVariant.WeightUnit,
				InventoryAvailable:                 updatedVariant.InventoryAvailable,
				InventoryCommitted:                 updatedVariant.InventoryCommitted,
				InventoryOnHand:                    updatedVariant.InventoryOnHand,
				InventoryUnavailableDamaged:        updatedVariant.InventoryUnavailableDamaged,
				InventoryUnavailableQualityControl: updatedVariant.InventoryUnavailableQualityControl,
				InventoryUnavailableSafetyStock:    updatedVariant.InventoryUnavailableSafetyStock,
				InventoryUnavailableOther:          updatedVariant.InventoryUnavailableOther,
			}
			// Create new SelectedOptions for the new variant
			for _, updatedSelectedOption := range updatedVariant.SelectedOptions {
				newSelectedOption := SelectedOption{
					VariantID: updatedSelectedOption.VariantID,
					Name:      updatedSelectedOption.Name,
					Value:     updatedSelectedOption.Value,
				}
				newVariant.SelectedOptions = append(newVariant.SelectedOptions, newSelectedOption)
			}
			existingProduct.Variants = append(existingProduct.Variants, newVariant)
		}
	}
	if err := pd.db.Save(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}

func (pd *ProductModel) DeleteProduct(id int) error {
	var existingProduct Product
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").First(&existingProduct, uint(id)).Error; err != nil {
		return errors.New("product not found")
	}

	// Delete the product and associated data
	if err := pd.db.Delete(&existingProduct).Error; err != nil {
		return errors.New("error deleting product: " + err.Error())
	}
	if err := pd.db.Delete(&existingProduct.Options).Error; err != nil {
		return errors.New("error deleting product's options: " + err.Error())
	}
	for _, variant := range existingProduct.Variants {
		if err := pd.deleteSelectedOptionsForVariant(variant.ID); err != nil {
			return errors.New("error deleting variant's seclected_Options: " + err.Error())
		}
	}
	if err := pd.db.Delete(&existingProduct.Variants).Error; err != nil {
		return errors.New("error deleting products's variants: " + err.Error())
	}
	return nil
}

func (pd *ProductModel) DeleteProductVariant(id int, vid int) error {
	var existingProduct Product
	if err := pd.db.Preload("Options").Preload("Variants").Preload("Variants.SelectedOptions").First(&existingProduct, uint(id)).Error; err != nil {
		return errors.New("product not found")
	}
	if len(existingProduct.Variants) == 1 {
		return errors.New("has only one varient")
	}
	// Find the index of the variant with the given ID
	variantIndex := -1
	for i, variant := range existingProduct.Variants {
		if variant.ID == uint(vid) {
			variantIndex = i
			break
		}
	}
	// Check if the variant with the given ID exists
	if variantIndex == -1 {
		return errors.New("variant not found")
	}
	// Delete the variant and its associated selected options
	pd.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existingProduct.Variants[variantIndex]).Error; err != nil {
			return err
		}

		// Delete the selected options associated with the variant
		if err := tx.Delete(&existingProduct.Variants[variantIndex].SelectedOptions).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (pd *ProductModel) GetPaginatedActiveProducts(page, limit int) ([]Product, error) {
	var offset int
	if page > 1 {
		offset = (page - 1) * limit
	}

	var existingProducts []Product
	if err := pd.db.Preload("Options").
		Preload("Variants").
		Preload("Variants.SelectedOptions").
		Offset(offset).
		Limit(limit).
		Find(&existingProducts, "status = ?", true).Error; err != nil {
		return nil, err
	}
	return existingProducts, nil
}

func (pd *ProductModel) deleteSelectedOptionsForVariant(variantID uint) error {
	// Find the variant
	var variant Variant
	if err := pd.db.First(&variant, "id = ?", variantID).Error; err != nil {
		return fmt.Errorf("error finding variant: %v", err)
	}

	// Find the associated SelectedOptions
	var selectedOptions []SelectedOption
	if err := pd.db.Where("variant_id = ?", variantID).Find(&selectedOptions).Error; err != nil {
		return fmt.Errorf("error finding selected options: %v", err)
	}

	if len(selectedOptions) == 0 {
		return nil
	}

	// Delete the associated SelectedOptions
	if err := pd.db.Delete(&selectedOptions).Error; err != nil {
		return fmt.Errorf("error deleting selected options: %v", err)
	}

	return nil
}
