package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID                        uint                  `gorm:"column:id;primaryKey"`
	Email                     string                `gorm:"column:email;index"`
	CreatedAt                 string                `gorm:"column:createdAt"`
	UpdatedAt                 string                `gorm:"column:updatedAt"`
	FirstName                 string                `gorm:"column:firstName"`
	LastName                  string                `gorm:"column:lastName"`
	OrdersCount               int                   `gorm:"column:ordersCount"`
	State                     string                `gorm:"column:state"`
	TotalSpent                float64               `gorm:"column:totalSpent"`
	LastOrderId               uint                  `gorm:"column:lastOrderId"`
	Note                      string                `gorm:"column:note"`
	VerifiedEmail             bool                  `gorm:"column:verifiedEmail"`
	MultipassIdentifier       bool                  `gorm:"column:multipassIdentifier"`
	TaxExempt                 bool                  `gorm:"column:taxExempt"`
	Tags                      bool                  `gorm:"column:tags"`
	LastOrderName             string                `gorm:"column:lastOrderName"`
	Currency                  string                `gorm:"column:currency"`
	Phone                     string                `gorm:"column:phone"`
	TaxExemptions             string                `gorm:"foreignKey:taxExemptions"`
	Password                  string                `gorm:"column:password;omitempty"`
	DisplayName               string                `gorm:"column:display_name"`
	Locale                    string                `gorm:"column:locale"`
	Age                       int                   `gorm:"column:age"`
	CustomerCreatTime         string                `gorm:"column:customerCreateTime"`
	LifetimeDuration          string                `gorm:"column:lifetimeDuration"`
	LastOrderCreatedAt        string                `gorm:"column:lastOrderCreatedAt"`
	AcceptsMarketing          bool                  `gorm:"column:acceptsMarketing"`
	AcceptsMarketingUpdatedAt string                `gorm:"column:acceptsMarketingUpdatedAt"`
	MarketingOptInLevel       string                `gorm:"column:marketingOptInLevel"`
	DeliveryAddresses         []DeliveryAddress     `gorm:"foreignKey:CustomerID"`
	EmailMarketingConsent     EmailMarketingConsent `gorm:"embedded"`
	SmsMarketingConsent       SmsMarketingConsent   `gorm:"embedded"`
	Metafield                 Metafield             `gorm:"embedded"`
}

type EmailMarketingConsent struct {
	State            string `gorm:"column:state"`
	OptInLevel       string `gorm:"column:optInLevel"`
	ConsentUpdatedAt string `gorm:"column:consentUpdatedAt"`
}

type Metafield struct {
	Key       string `gorm:"column:key"`
	Namespace string `gorm:"column:namespace"`
	Value     string `gorm:"column:value"`
	Type      string `gorm:"column:type"`
}

type SmsMarketingConsent struct {
	State                string `gorm:"column:state"`
	OptInLevel           string `gorm:"column:optInLevel"`
	ConsentUpdatedAt     string `gorm:"column:consentUpdatedAt"`
	ConsentCollectedFrom string `gorm:"column:consentCollectedFrom"`
}

type DeliveryAddress struct {
	ID           uint   `gorm:"column:id;primaryKey"`
	CustomerID   uint   `gorm:"column:customer_id;index"`
	FirstName    string `gorm:"column:firstName"`
	LastName     string `gorm:"column:lastName"`
	Company      string `gorm:"column:company"`
	Address1     string `gorm:"column:address1"`
	Address2     string `gorm:"column:address2"`
	Apartment    string `gorm:"column:apartment"`
	City         string `gorm:"column:city"`
	Province     string `gorm:"column:province"`
	Country      string `gorm:"column:country"`
	Zip          int    `gorm:"column:zip"`
	Phone        int    `gorm:"column:phone"`
	ProvinceCode string `gorm:"column:provinceCode"`
	CountryCode  string `gorm:"column:countryCode"`
	CountryName  string `gorm:"column:countryName"`
	Default      bool   `gorm:"column:default"`
}

type CustomerSecrets struct {
	ID         uint   `gorm:"primaryKey" json:"id" `
	CustomerID uint   `gorm:"column:customer_id"`
	Email      string `gorm:"column:email"`
	Secret     string `gorm:"column:secret"`
}

type CustomerToken struct {
	ID                uint      `gorm:"primaryKey" json:"id" `
	Email             string    `gorm:"column:email;primaryKey"`
	Token             string    `gorm:"column:token"`
	LastUsed          time.Time `gorm:"column:last_used"`
	TokenExpiry       time.Time `gorm:"column:token_expiry"`
	IsActive          bool      `gorm:"column:is_active"`
	IPAddresses       string    `gorm:"column:ip_addresses"`
	CustomerAgent     string    `gorm:"column:customer_agent"`
	DeviceInformation string    `gorm:"column:device_information"`
	RevocationReason  string    `gorm:"column:revocation_reason"`
}

type CustomerModel struct {
	db *gorm.DB
}

func NewCustomerModel(db *gorm.DB) *CustomerModel {
	return &CustomerModel{db: db}
}

func (ur *CustomerModel) Save(customer *Customer) error {
	return ur.db.Create(customer).Error
}

func (ur *CustomerModel) SaveCustomerSecret(CustomerSecret *CustomerSecrets) error {
	return ur.db.Create(CustomerSecret).Error
}

func GetCustomerSecretKeyByEmail(db *gorm.DB, email string) (string, error) {
	var customerSecrets CustomerSecrets
	err := db.Model(&CustomerSecrets{}).Select("secret").Where("email = ?", email).First(&customerSecrets).Error
	if err != nil {
		return "", err
	}
	return customerSecrets.Secret, nil
}

func GetCustomerHashedPasswordByEmail(db *gorm.DB, email string) (string, error) {
	var CustomerHashedPasswor Customer
	err := db.Model(&Customer{}).Select("Password").Where("email = ?", email).First(&CustomerHashedPasswor).Error
	if err != nil {
		return "", err
	}
	return CustomerHashedPasswor.Password, nil
}

func CustomerExistByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&Customer{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func SaveToken(db *gorm.DB, customerToken *CustomerToken) error {
	return db.Create(customerToken).Error
}

func CheckEmailExist(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&CustomerToken{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func GetTokenByEmail(db *gorm.DB, email string) (string, error) {
	var customerToken CustomerToken
	err := db.Model(&CustomerToken{}).Select("token").Where("email = ?", email).First(&customerToken).Error
	if err != nil {
		return "", err
	}
	return customerToken.Token, nil
}
func UpdateToken(db *gorm.DB, customerToken *CustomerToken) error {
	updatedValues := map[string]interface{}{
		"email":              customerToken.Email,
		"token":              customerToken.Token,
		"last_used":          customerToken.LastUsed,
		"token_expiry":       customerToken.TokenExpiry,
		"is_active":          customerToken.IsActive,
		"ip_addresses":       customerToken.IPAddresses,
		"customer_agent":     customerToken.CustomerAgent,
		"device_information": customerToken.DeviceInformation,
		"revocation_reason":  customerToken.RevocationReason,
	}
	return db.Model(&CustomerToken{}).Where("email = ?", customerToken.Email).Updates(updatedValues).Error
}

func (ad *CustomerModel) GetCustomerID(email string) (uint, error) {
	var customer Customer
	if err := ad.db.Model(&Customer{}).Where("email = ?", email).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("customer not found")
		}
		return 0, fmt.Errorf("error fetching customer: %v", err)
	}
	return customer.ID, nil
}
