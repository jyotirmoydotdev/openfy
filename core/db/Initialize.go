package db

import (
	"fmt"

	"github.com/jyotirmoydotdev/openfy/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var productdbInstance *gorm.DB

func GetDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	db, err := gorm.Open(sqlite.Open("./db/databaseUserAdmin.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	dbInstance = db
	return db, nil
}

func GetProductDB() (*gorm.DB, error) {
	if productdbInstance != nil {
		return productdbInstance, nil
	}
	db, err := gorm.Open(sqlite.Open("./db/databaseProduct.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	productdbInstance = db
	return db, nil
}

func InitializeDatabases() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSecrets{},
		&models.DeliveryAddress{},
		&models.ShopDetail{},
		&models.Admin{},
		&models.AdminSecrets{},
		&models.Counter{},
		&models.UserToken{},
	)
	if err != nil {
		return fmt.Errorf("error auto migrating models: %v", err)
	}
	err = preloadData(db)
	if err != nil {
		return fmt.Errorf("error preloading data: %v", err)
	}

	productdb, err := GetProductDB()
	if err != nil {
		return err
	}
	err = productdb.AutoMigrate(
		&models.Product{},
		&models.Option{},
		&models.Variant{},
		&models.SelectedOption{},
		&models.Counter{},
	)
	if err != nil {
		return fmt.Errorf("error auto migrating models: %v", err)
	}
	return nil
}

func preloadData(db *gorm.DB) error {
	err := db.Create(&models.Counter{
		Name:  "admin",
		Count: 1,
	}).Error
	if err != nil {
		return err
	}
	err = db.Create(&models.Counter{
		Name:  "user",
		Count: 1,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
