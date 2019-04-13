// work with table keys
package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

const TableKeys = "keys"

type Keys struct {
	Id     int `gorm:"primary_key"`
	Key    string
	Name   string
	Ip     string
	Status int8
}

func (obj *Keys) checkKeyIsSet() error {
	if obj.Key == "" {
		return fmt.Errorf("error: no key to greate row")
	}
	return nil
}

func (obj *Keys) Create(db *gorm.DB) error {
	if err := obj.checkKeyIsSet(); err != nil {
		return err
	}
	if err := db.Table(TableKeys).
		Create(obj).Error; err != nil {
		return err
	} else {
		return nil
	}
}

func (obj *Keys) LoadByKey(db *gorm.DB) error {
	if err := obj.checkKeyIsSet(); err != nil {
		return err
	}
	if err := db.Table(TableKeys).
		Where("key = ?", obj.Key).First(obj).Error; err != nil {
		return err
	}
	return nil
}

// GetKeysByStatus get keys by status
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
