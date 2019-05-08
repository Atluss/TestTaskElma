// Package config this package for create DB conn and Route from config file.
package config

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

// NewApiSetup return server config with gorm DB and router
func NewApiSetup(settings string) *Setup {

	cnf, err := Config(settings)
	v1.FailOnError(err, "error config file")

	set, err := newSetup(cnf)
	v1.FailOnError(err, "error setup")

	return set
}

func newSetup(cnf *config) (*Setup, error) {

	set := Setup{}

	if err := cnf.validate(); err != nil {
		return &set, err
	}

	set.Config = cnf
	set.Route = mux.NewRouter().StrictSlash(true)

	if err := set.getDB(); err != nil {
		return &set, err
	}

	return &set, nil
}

// Setup main setup api struct
type Setup struct {
	Config *config     // api setting
	Route  *mux.Router // mux frontend
	Gorm   *gorm.DB
}

// getDB setup gorm connection for DB
func (obj *Setup) getDB() error {

	connectQuery := fmt.Sprintf(obj.Config.Gorm.ConnPattern, obj.Config.Gorm.Type, obj.Config.Gorm.User,
		obj.Config.Gorm.Password, obj.Config.Gorm.Host, obj.Config.Gorm.Port, obj.Config.Gorm.Database)

	db, err := gorm.Open(obj.Config.Gorm.Type, connectQuery)
	if err != nil {
		log.Println(err)
		return err
	}
	obj.Gorm = db

	return nil
}
