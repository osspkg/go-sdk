package main

import (
	"fmt"

	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/log"
)

type (
	//Simple model
	Simple struct{}
	//Config model
	Config1 struct {
		Env string `yaml:"env"`
	}
	Config2 struct {
		Env string `yaml:"env"`
	}
)

// NewSimple init Simple
func NewSimple(c1 Config1, c2 Config2) *Simple {
	fmt.Println("--> call NewSimple")
	fmt.Println("--> Config1.ENV=" + c1.Env)
	fmt.Println("--> Config2.ENV=" + c2.Env)
	return &Simple{}
}

// Up  method for start Simple in DI container
func (s *Simple) Up(_ app.Context) error {
	fmt.Println("--> call *Simple.Up")
	return nil
}

// Down  method for stop Simple in DI container
func (s *Simple) Down(_ app.Context) error {
	fmt.Println("--> call *Simple.Down")
	return nil
}

func main() {
	app.New().
		Logger(log.Default()).
		ConfigFile(
			"./config.yaml",
			Config1{},
		).
		Modules(
			Config2{Env: "prod"},
			NewSimple,
		).
		Run()
}
