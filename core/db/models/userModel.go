package models

import "gorm.io/gorm"

type User struct {
	ID                string            `gorm:"column:id;primaryKey"`
	Password          string            `gorm:"column:password;omitempty"`
	Email             string            `gorm:"column:email;index"`
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
	Email  string `gorm:"column:email"`
	Secret string `gorm:"column:secret"`
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
