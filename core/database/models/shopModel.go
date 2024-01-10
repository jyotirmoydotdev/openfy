package models

import "gorm.io/gorm"

type ShopDetail struct {
	gorm.Model
	ID                      uint   `gorm:"column:id"`
	Name                    string `gorm:"column:name"`
	Email                   string `gorm:"column:email"`
	CurrencyCode            string `gorm:"column:currencyCode"`
	IanaTimezone            string `gorm:"column:ianaTimezone"`
	TimezoneOffsetMinutes   int    `gorm:"column:timezoneOffsetMinutes"`
	TimezoneAbbreviation    string `gorm:"column:timezoneAbbreviation"`
	ShopDomain              string `gorm:"column:shopDomain"`
	CreatedAt               string `gorm:"column:createdAt"`
	OrderNumberFormatPrefix string `gorm:"column:orderNumberFormatPrefix"`
	OrderNumberFormatSuffix string `gorm:"column:orderNumberFormatSuffix"`
	UnitSystem              string `gorm:"column:unitSystem"`
	WightUnit               string `gorm:"column:weightUnit"`
}
