// Package dataKeys work with table keys
package dataKeys

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

// keys status
const (
	// TableKeys table name
	TableKeys = "keys"
	// KeyEmpty status empty
	KeyEmpty = 0
	// KeyActive status accept
	KeyActive = 1
	// KeyBlocked status blocked
	KeyBlocked = 2
	// KeyOnline status online no for DB
	KeyOnline = 100
	// KeyOffline status go offline no for DB
	KeyOffline = 200
)

// Keys row in database
type Keys struct {
	Id     int `gorm:"primary_key"`
	Key    string
	Name   string
	Ip     string
	Status int
}

func (obj *Keys) checkKeyIsSet() error {
	if obj.Key == "" {
		return fmt.Errorf("error: no key to create row")
	}
	return nil
}

// Create key in DB
func (obj *Keys) Create(db *gorm.DB) error {
	if err := obj.checkKeyIsSet(); err != nil {
		return err
	}
	if err := db.Table(TableKeys).
		Create(obj).Error; err != nil {
		return err
	}
	return nil
}

// LoadByKey load key from DB
func (obj *Keys) LoadByKey(db *gorm.DB) error {
	if err := obj.checkKeyIsSet(); err != nil {
		return err
	}
	return db.Table(TableKeys).Where("key = ?", obj.Key).First(obj).Error
}

// Update key in DB
func (obj *Keys) Update(db *gorm.DB) error {
	if err := obj.checkKeyIsSet(); err != nil {
		return err
	}

	if err := db.Table(TableKeys).Where("key = ?", obj.Key).Update(obj).Error; err != nil {
		return fmt.Errorf("error to update key: %s", err)
	}
	return nil
}

// GetKeysByStatus get keys by status if status more than 10 returns all keys
func GetKeysByStatus(status int, db *gorm.DB) (keys []Keys, err error) {
	tx := db.Table(TableKeys)
	if status < 3 {
		tx = tx.Where("status = ?", status)
	}
	if err := tx.Scan(&keys).Error; err != nil {
		log.Printf("error: get keys by status %d, err: %s", status, err)
		return keys, err
	}
	return keys, err
}
