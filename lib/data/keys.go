// work with table keys
package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

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
	gorm.Model
	Key    string
	Status int8
	Info   Info
}
