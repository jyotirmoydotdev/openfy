package database

import (
	"fmt"

	"github.com/jyotirmoydotdev/openfy/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var customerInstance *gorm.DB
var staffMember *gorm.DB
var productInstance *gorm.DB

func GetCustomerDB() (*gorm.DB, error) {
	if customerInstance != nil {
		return customerInstance, nil
	}
	db, err := gorm.Open(sqlite.Open("./database/databaseCustomer.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening databaseCustomer: %v", err)
	}
	customerInstance = db
	return db, nil
}

func GetStaffMemberDB() (*gorm.DB, error) {
	if staffMember != nil {
		return staffMember, nil
	}
	db, err := gorm.Open(sqlite.Open("./database/databaseStaffMember.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening databaseStaffMember: %v", err)
	}
	staffMember = db
	return db, nil
}

func GetProductDB() (*gorm.DB, error) {
	if productInstance != nil {
		return productInstance, nil
	}
	db, err := gorm.Open(sqlite.Open("./database/databaseProduct.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening databaseProduct: %v", err)
	}
	productInstance = db
	return db, nil
}

func InitializeDatabases() error {
	customerDB, err := GetCustomerDB()
	if err != nil {
		return err
	}
	err = customerDB.AutoMigrate(
		&models.Customer{},
		&models.CustomerSecrets{},
		&models.DeliveryAddress{},
		&models.ShopDetail{},
		&models.CustomerToken{},
	)
	if err != nil {
		return fmt.Errorf("error auto migrating customerDB: %v", err)
	}

	staffMemberDB, err := GetStaffMemberDB()
	if err != nil {
		return err
	}

	err = staffMemberDB.AutoMigrate(
		&models.StaffMember{},
		&models.StaffMemberSecrets{},
	)
	if err != nil {
		return fmt.Errorf("error auto migrating staffMemberDB: %v", err)
	}

	productDB, err := GetProductDB()
	if err != nil {
		return err
	}
	err = productDB.AutoMigrate(
		&models.Product{},
		&models.Option{},
		&models.Variant{},
		&models.SelectedOption{},
	)
	if err != nil {
		return fmt.Errorf("error auto migrating GetProductDB: %v", err)
	}
	return nil
}
