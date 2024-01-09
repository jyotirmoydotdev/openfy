package models

import (
	"fmt"

	"gorm.io/gorm"
)

type StaffMember struct {
	ID           uint     `gorm:"column:id;primaryKey"`
	Customername string   `gorm:"column:customername;index"`
	Password     string   `gorm:"column:password"`
	Email        string   `gorm:"column:email"`
	Name         string   `gorm:"column:name"`
	FirstName    string   `gorm:"column:firstName"`
	LastName     string   `gorm:"column:lastName"`
	AccountOwner bool     `gorm:"column:accountOwner"`
	Locale       string   `gorm:"column:locale"`
	Permissions  []string `gorm:"column:permissions;type:text[]"`
}

type AdminSecrets struct {
	AdminID      uint   `gorm:"column:admin_id"`
	Customername string `gorm:"column:customername"`
	Secret       string `gorm:"column:secret"`
}

type AdminModel struct {
	db *gorm.DB
}

func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{db: db}
}

func (ad *AdminModel) Save(admin *StaffMember) error {
	return ad.db.Create(admin).Error
}
func (ad *AdminModel) SaveAdminSecret(adminSecret *AdminSecrets) error {
	return ad.db.Create(adminSecret).Error
}
func AdminExistByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func AdminExistByCustomername(db *gorm.DB, customername string) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Where("customername = ?", customername).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func CheckAdminTableIsEmpty(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
func GetAdminHashedPasswordByCustomername(db *gorm.DB, customername string) (string, error) {
	var admin StaffMember
	if err := db.Model(&StaffMember{}).Where("customername = ?", customername).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error fetching customer: %v", err)
	}

	return admin.Password, nil
}
func GetSecretKeyByCustomername(db *gorm.DB, customername string) (string, error) {
	var admin AdminSecrets
	if err := db.Model(&AdminSecrets{}).Where("customername = ?", customername).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error fetching customer: %v", err)
	}

	return admin.Secret, nil
}
func (ad *AdminModel) GetAdminID(customername string) (uint, error) {
	var admin StaffMember
	if err := ad.db.Model(&StaffMember{}).Where("customername = ?", customername).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("admin not found")
		}
		return 0, fmt.Errorf("error fetching customer: %v", err)
	}
	return admin.ID, nil
}
