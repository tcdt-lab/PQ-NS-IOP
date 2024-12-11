package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadYaml(t *testing.T) {
	t.Log("Testing ReadYaml")
	c, err := ReadYaml()
	if err != nil {
		t.Errorf("Error in ReadYaml: %v", err)
	}
	assert.NoError(t, err, "Error in ReadYaml")
	t.Log(c)
}
