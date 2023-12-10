package models

import "gorm.io/gorm"

type ShopDetail struct {
	gorm.Model
	Name                  string `gorm:"column:name"`
	Description           string `gorm:"column:description"`
	TimeZone              string `gorm:"column:time_zone"`
	TimeZoneOffsetMinutes int    `gorm:"column:time_zone_offset_minutes"`
	TimezoneOffsetMinutes string `gorm:"column:timezone_offset_minutes"`
	ShopDomain            string `gorm:"column:shop_domain"`
	CreatedAt             struct {
		CreatedAt string `gorm:"column:created_at"`
		Date      int    `gorm:"column:date"`
		Month     int    `gorm:"column:month"`
		Year      int    `gorm:"column:year"`
		Time      struct {
			Hour   int `gorm:"column:hour"`
			Minute int `gorm:"column:minute"`
			Second int `gorm:"column:second"`
		} `gorm:"embedded"`
	} `gorm:"embedded"`
	Location struct {
		Address string `gorm:"column:address"`
		City    string `gorm:"column:city"`
		State   string `gorm:"column:state"`
		ZipCode string `gorm:"column:zip_code"`
		Country string `gorm:"column:country"`
	} `gorm:"embedded"`
	Contact struct {
		Email       string `gorm:"column:email"`
		Phone       string `gorm:"column:phone"`
		WhatsApp    string `gorm:"column:whatsapp"`
		TwitterURL  string `gorm:"column:twitter_url"`
		FacebookURL string `gorm:"column:facebook_url"`
		YoutubeURL  string `gorm:"column:youtube_url"`
	} `gorm:"embedded"`
	Currency struct {
		Code   string `gorm:"column:code"`
		Symbol string `gorm:"column:symbol"`
	} `gorm:"embedded"`
	UintSystem struct {
		MerticSystem struct {
			Kilogram bool `gorm:"column:kilogram"`
			Gram     bool `gorm:"column:gram"`
		} `gorm:"embedded"`
		ImperialSystem struct {
			Pound bool `gorm:"column:pound"`
			Ounce bool `gorm:"column:ounce"`
		} `gorm:"embedded"`
	} `gorm:"embedded"`
}
