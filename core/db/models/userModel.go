package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string            `gorm:"column:id;primaryKey"`
	Password          string            `gorm:"column:password;omitempty"`
	Email             string            `gorm:"column:email;index"`
	FirstName         string            `gorm:"column:first_name"`
	LastName          string            `gorm:"column:last_name"`
	Phone             int               `gorm:"column:phone"`
	Age               int               `gorm:"column:age"`
	DeliveryAddresses []DeliveryAddress `gorm:"foreignKey:UserID"`
}

type DeliveryAddress struct {
	ID        string `gorm:"column:id;primaryKey"`
	UserID    string `gorm:"column:user_id;index"`
	Country   string `gorm:"column:country"`
	Address   string `gorm:"column:address"`
	Apartment string `gorm:"column:apartment"`
	City      string `gorm:"column:city"`
	State     string `gorm:"column:state"`
	PinCode   int    `gorm:"column:pincode"`
}

type UserSecrets struct {
	UserID string `gorm:"column:user_id"`
	Email  string `gorm:"column:email"`
	Secret string `gorm:"column:secret"`
}

type UserToken struct {
	Email             string    `gorm:"column:email;primaryKey"`
	Token             string    `gorm:"column:token"`
	LastUsed          time.Time `gorm:"column:last_used"`
	TokenExpiry       time.Time `gorm:"column:token_expiry"`
	IsActive          bool      `gorm:"column:is_active"`
	IPAddresses       string    `gorm:"column:ip_addresses"`
	UserAgent         string    `gorm:"column:user_agent"`
	DeviceInformation string    `gorm:"column:device_information"`
	RevocationReason  string    `gorm:"column:revocation_reason"`
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

func (ur *UserModel) Save(user *User) error {
	return ur.db.Create(user).Error
}

func (ur *UserModel) SaveUserSecret(UserSecret *UserSecrets) error {
	return ur.db.Create(UserSecret).Error
}

func GetUserSecretKeyByEmail(db *gorm.DB, email string) (string, error) {
	var userSecrets UserSecrets
	err := db.Model(&UserSecrets{}).Select("secret").Where("email = ?", email).First(&userSecrets).Error
	if err != nil {
		return "", err
	}
	return userSecrets.Secret, nil
}

func GetUserHashedPasswordByEmail(db *gorm.DB, email string) (string, error) {
	var UserHashedPasswor User
	err := db.Model(&User{}).Select("Password").Where("email = ?", email).First(&UserHashedPasswor).Error
	if err != nil {
		return "", err
	}
	return UserHashedPasswor.Password, nil
}

func UserExistByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func SaveToken(db *gorm.DB, userToken *UserToken) error {
	return db.Create(userToken).Error
}

func CheckEmailExist(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&UserToken{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func GetTokenByEmail(db *gorm.DB, email string) (string, error) {
	var userToken UserToken
	err := db.Model(&UserToken{}).Select("token").Where("email = ?", email).First(&userToken).Error
	if err != nil {
		return "", err
	}
	return userToken.Token, nil
}
func UpdateToken(db *gorm.DB, userToken *UserToken) error {
	updatedValues := map[string]interface{}{
		"email":              userToken.Email,
		"token":              userToken.Token,
		"last_used":          userToken.LastUsed,
		"token_expiry":       userToken.TokenExpiry,
		"is_active":          userToken.IsActive,
		"ip_addresses":       userToken.IPAddresses,
		"user_agent":         userToken.UserAgent,
		"device_information": userToken.DeviceInformation,
		"revocation_reason":  userToken.RevocationReason,
	}
	return db.Model(&UserToken{}).Where("email = ?", userToken.Email).Updates(updatedValues).Error
}
