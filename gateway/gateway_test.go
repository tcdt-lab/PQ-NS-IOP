package main

import (
	"gateway/config"
	"testing"
)

func TestBootLogic(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file")
	}
	BootLogic(c)
}
func TestPhaseOneExecute(t *testing.T) {
	PhaseOneExecute()
}
