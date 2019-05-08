package config

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"io/ioutil"
	"log"
	"os"
)

// Config load new config for API
func Config(path string) (*config, error) {

	conf := config{}

	if err := v1.CheckFileExist(path); err != nil {
		return &conf, err
	}

	conf.FilePath = path

	if err := conf.load(); err != nil {
		return &conf, err
	}

	return &conf, nil
}

type gormConfig struct {
	Type        string `json:"Type"`
	Host        string `json:"Host"`
	Port        string `json:"Port"`
	User        string `json:"User"`
	Password    string `json:"Password"`
	Database    string `json:"Database"`
	ConnPattern string `json:"ConnPattern"`
}

// config main
type config struct {
	Name     string     `json:"Name"`     // API name
	Version  string     `json:"Version"`  // API version
	Host     string     `json:"Host"`     // API host
	Port     string     `json:"Port"`     // API port
	FilePath string     `json:"FilePath"` // path to Json settings file
	Gorm     gormConfig `json:"Gorm"`
}

// load all settings
func (obj *config) load() error {

	jsonSet, err := os.Open(obj.FilePath)

	defer func() {
		// defer and handle close error
		v1.LogOnError(jsonSet.Close(), "warning: Can't close json settings file.")
	}()

	if !v1.LogOnError(err, "Can't open config file") {
		return err
	}

	bytesVal, _ := ioutil.ReadAll(jsonSet)
	err = json.Unmarshal(bytesVal, &obj)

	if !v1.LogOnError(err, "Can't unmarshal json file") {
		return err
	}

	return obj.validate()
}

// validate it
func (obj *config) validate() error {

	if obj.Name == "" {
		return fmt.Errorf("config miss name")
	}

	if obj.Version == "" {
		return fmt.Errorf("config miss version")
	}

	if obj.Host == "" {
		return fmt.Errorf("config miss address")
	}

	if obj.Port == "" {
		return fmt.Errorf("config miss port")
	}

	if obj.Gorm.Type == "" {
		return fmt.Errorf("config miss gorm type")
	}

	if obj.Gorm.Host == "" {
		return fmt.Errorf("config miss gorm host")
	}

	if obj.Gorm.Port == "" {
		return fmt.Errorf("config miss gorm port")
	}

	if obj.Gorm.User == "" {
		return fmt.Errorf("config miss gorm user")
	}

	if obj.Gorm.Password == "" {
		return fmt.Errorf("config miss gorm password")
	}

	if obj.Gorm.ConnPattern == "" {
		return fmt.Errorf("config miss gorm connPatter")
	}

	if obj.Gorm.Database == "" {
		return fmt.Errorf("config miss gorm database")
	}

	return nil
}

// Print config
func (obj *config) Print() {
	log.Printf("Name: %s", obj.Name)
	log.Printf("Version: %s", obj.Version)
	log.Printf("Host: %s", obj.Host)
	log.Printf("Port: %s", obj.Port)
}
