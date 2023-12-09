package database

type ShopDetail struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	TimeZone              string `json:"timeZone"`
	TimeZoneOffsetMinutes int    `json:"timeZoneOffsetMinutes"`
	TimezoneOffsetMinutes string `json:"timezoneOffsetMinutes"`
	ShopDomain            string `json:"shopDomain"`
	CreatedAt             struct {
		CreatedAt string `json:"createdAt"`
		Date      int    `json:"date"`
		Month     int    `json:"month"`
		Year      int    `json:"year"`
		Time      struct {
			Hour   int `json:"hour"`
			Minute int `json:"minute"`
			Second int `json:"second"`
		} `json:"time"`
	} `json:"createdAt"`
	Location struct {
		Address string `json:"address"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode string `json:"zipCode"`
		Country string `json:"country"`
	} `json:"location"`
	Contact struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		WhatsApp    string `json:"whatsApp"`
		TwitterURL  string `json:"twitterURL"`
		FacebookURL string `json:"facebookURL"`
		YoutubeURL  string `json:"youtubeURL"`
	} `json:"contact"`
	Currency struct {
		Code   string `json:"code"`
		Symbol string `json:"symbol"`
	} `json:"currency"`
	UintSystem struct {
		MerticSystem struct {
			Kilogram bool `json:"kilogram"`
			Gram     bool `json:"gram"`
		} `json:"merticSystem"`
		ImperialSystem struct {
			Pound bool `json:"pound"`
			Ounce bool `json:"ounce"`
		} `json:"imperialSystem"`
	} `json:"uintSystem"`
}

var ShopDetails ShopDetail
