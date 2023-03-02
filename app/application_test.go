package app_test

import "testing"

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

func TestUnit_Invoke(t *testing.T) {

}
