package datamodels

type ProductRequest struct {
	Title           string `json:"title"`
	DescriptionHtml string `json:"descriptionHtml"`
	Handle          string `json:"handle"`
	SEO             struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"seo"`
	Tags                   []string `json:"tags"`
	TemplateSuffix         string   `json:"templateSuffix"`
	GiftCardTemplateSuffix string   `json:"giftCardTemplateSuffix"`
	Vendor                 string   `json:"vendor"`
	ProductCategory        string   `json:"productCategory"`
	ProductType            string   `json:"productType"`
	Publications           []struct {
		PublicationID string `json:"publicationID"`
	} `json:"publications"`
	GiftCard          bool     `json:"giftCard"`
	CollectionsToJoin []string `json:"CollectionsToJoin"`
	Workflow          string   `json:"workflow"`
	Metafields        []string `json:"metafields"`
	Media             []struct {
		MediaContentType string `json:"mediaContentType"`
		OrginalSource    string `json:"orginalSource"`
	} `json:"media"`
}

type ProductResponse struct {
	ID             uint              `json:"id"`
	Title          string            `json:"title"`
	BodyHTML       string            `json:"bodyHtml"`
	Vendor         string            `json:"vendor"`
	ProductType    string            `json:"productType"`
	CreatedAt      string            `json:"createAt"`
	Handle         string            `json:"handle"`
	UpdatedAt      string            `json:"updatedAt"`
	PublishedAt    string            `json:"publishedAt"`
	TemplateSuffix string            `json:"templateSuffix"`
	PublishedScope string            `json:"publishScope"`
	Tags           string            `json:"tags"`
	Status         string            `json:"status"`
	Variants       []VariantResponse `json:"variants"`
	Options        []Option          `json:"options"`
	Images         []Image           `json:"images"`
	Image          Image             `json:"image"`
}

type VariantResponse struct {
	ID                  uint               `json:"id"`
	ProductID           uint               `json:"productID"`
	Title               string             `json:"title"`
	Price               string             `json:"price"`
	SKU                 string             `json:"sku"`
	Position            uint               `json:"position"`
	InventoryPolicy     string             `json:"inventoryPolicy"`
	CompareAtPrice      string             `json:"compareAtPrice"`
	FullfillmentService string             `json:"fullfillmentService"`
	InventoryManagement string             `json:"inventoryManagement"`
	Option1             string             `json:"Option1"`
	Option2             string             `json:"Option2"`
	Option3             string             `json:"Option3"`
	CreatedAt           string             `json:"createAt"`
	UpdatedAt           string             `json:"updatedAt"`
	Taxable             bool               `json:"taxable"`
	Barcode             string             `json:"barcode"`
	Grams               uint               `json:"garms"`
	ImageID             uint               `json:"imageID"`
	Weight              uint               `json:"weight"`
	WeightUint          string             `json:"weightUint"`
	InventoryItemID     uint               `json:"inventoryItemID"`
	InventoryQuantity   uint               `json:"inventoryQuantity"`
	PresentmentPrices   []PresentmentPrice `json:"presentmentPrices"`
}

type Option struct {
	ID        uint     `json:"id"`
	ProductID uint     `json:"productID"`
	Name      string   `json:"name"`
	Values    []string `json:"values"`
}

type Image struct {
	ID         uint   `json:"id"`
	ProductID  uint   `json:"productID"`
	Position   uint   `json:"position"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Alt        string `json:"alt"`
	Width      uint   `json:"width"`
	Height     uint   `json:"height"`
	Src        string `json:"src"`
	VariantIDs []uint `json:"variantIDs"`
}

type PresentmentPrice struct {
	Price          string `json:"price"`
	CompareAtPrice string `json:"compareAtPrice"`
	CurrencyCode   string `json:"currencyCode"`
}
