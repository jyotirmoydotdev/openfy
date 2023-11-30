package database

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
var ProductIDCounter int
