package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Counter struct {
	ID    int    `gorm:"column:ID;primaryKey"`
	Name  string `gorm:"column:name"`
	Count int    `gorm:"column:count"`
}

type CounterModel struct {
	db *gorm.DB
}

func NewCounterModel(db *gorm.DB) *CounterModel {
	return &CounterModel{db: db}
}
func (cu *CounterModel) Save(counter *Counter) error {
	return cu.db.Create(counter).Error
}
func Save(dbInstance *gorm.DB, counter *Counter) error {
	return dbInstance.Create(counter).Error
}
func Increment(dbInstance *gorm.DB, name string) error {
	var counter Counter
	if err := dbInstance.Where("name = ?", name).First(&counter).Error; err != nil {
		return fmt.Errorf("error finding counter: %v", err)
	}
	if counter.ID == 0 {
		counter = Counter{Name: name, Count: 0}
		if err := dbInstance.Create(&counter).Error; err != nil {
			return fmt.Errorf("error creating counter: %v", err)
		}
	}

	// Increment the counter
	counter.Count++

	// Update the counter in the database
	if err := dbInstance.Save(&counter).Error; err != nil {
		return fmt.Errorf("error updating counter: %v", err)
	}

	return nil
}
func GetCountOf(dbInstance *gorm.DB, name string) (int, error) {
	var counter Counter
	if err := dbInstance.Model(&Counter{}).Where("name = ?", name).First(&counter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("Counter not found")
		}
		return 0, fmt.Errorf("error fetching Counter: %v", err)
	}
	return counter.Count, nil
}
