// work with table keys
package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

const TableKeys = "keys"

type Info struct {
	Name string `json:"Name"`
	Ip   string `json:"Ip"`
}

func (obj *Info) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.([]byte); ok {
		if err := json.Unmarshal(s, &obj); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error convert value to string")
}

func (obj Info) Value() (driver.Value, error) {
	if data, err := json.Marshal(obj); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

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
