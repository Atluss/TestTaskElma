package config

import (
	"github.com/Atluss/TestTaskElma/lib"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {

	path := "settings.json"

	cnf, err := Config(path)
	lib.FailOnError(err, "Test error")
	log.Printf("%+v", cnf)
}

func TestSetup(t *testing.T) {

	path := "settings.json"
	set := NewApiSetup(path)

	log.Printf("%+v", set)
}
