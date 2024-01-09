package models

import (
	"fmt"

	"gorm.io/gorm"
)

type StaffMember struct {
	ID           uint     `gorm:"column:id;primaryKey"`
	Username     string   `gorm:"column:username;index"`
	Password     string   `gorm:"column:password"`
	Email        string   `gorm:"column:email"`
	Name         string   `gorm:"column:name"`
	FirstName    string   `gorm:"column:firstName"`
	LastName     string   `gorm:"column:lastName"`
	AccountOwner bool     `gorm:"column:accountOwner"`
	Locale       string   `gorm:"column:locale"`
	Permissions  []string `gorm:"column:permissions;type:text[]"`
}

type StaffMemberSecrets struct {
	StaffMemberID uint   `gorm:"column:staffMember_id"`
	Username      string `gorm:"column:username"`
	Secret        string `gorm:"column:secret"`
}

type StaffMemberModel struct {
	db *gorm.DB
}

func NewStaffMemberModel(db *gorm.DB) *StaffMemberModel {
	return &StaffMemberModel{db: db}
}

func (ad *StaffMemberModel) Save(staffMember *StaffMember) error {
	return ad.db.Create(staffMember).Error
}
func (ad *StaffMemberModel) SaveStaffMemberSecret(staffMemberSecret *StaffMemberSecrets) error {
	return ad.db.Create(staffMemberSecret).Error
}
func StaffMemberExistByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func StaffMemberExistByUsername(db *gorm.DB, username string) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func CheckStaffMemberTableIsEmpty(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&StaffMember{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
func GetStaffMemberHashedPasswordByUsername(db *gorm.DB, username string) (string, error) {
	var staffMember StaffMember
	if err := db.Model(&StaffMember{}).Where("username = ?", username).First(&staffMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("staffMember not found")
		}
		return "", fmt.Errorf("error fetching customer: %v", err)
	}

	return staffMember.Password, nil
}
func GetSecretKeyByUsername(db *gorm.DB, username string) (string, error) {
	var staffMember StaffMemberSecrets
	if err := db.Model(&StaffMemberSecrets{}).Where("username = ?", username).First(&staffMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("staffMember not found")
		}
		return "", fmt.Errorf("error fetching customer: %v", err)
	}

	return staffMember.Secret, nil
}
func (ad *StaffMemberModel) GetStaffMemberID(username string) (uint, error) {
	var staffMember StaffMember
	if err := ad.db.Model(&StaffMember{}).Where("username = ?", username).First(&staffMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("staffMember not found")
		}
		return 0, fmt.Errorf("error fetching customer: %v", err)
	}
	return staffMember.ID, nil
}
