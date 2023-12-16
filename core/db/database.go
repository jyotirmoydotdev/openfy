package db

import (
	"fmt"

	"github.com/jyotirmoydotdev/openfy/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func GetDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	dbInstance = db
	return db, nil
}

func InitializeDatabases() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.User{}, &models.UserSecrets{}, &models.DeliveryAddress{}, &models.ShopDetail{}, &models.Admin{}, &models.AdminSecrets{})
	if err != nil {
		return fmt.Errorf("error auto migrating models: %v", err)
	}
	return nil
}