package config

import (
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {

	path := "settings.json"

	cnf, err := Config(path)
	v1.FailOnError(err, "Test error")
	log.Printf("%+v", cnf)
}

func TestSetup(t *testing.T) {

	path := "settings.json"
	set := NewApiSetup(path)

	log.Printf("%+v", set)
}
