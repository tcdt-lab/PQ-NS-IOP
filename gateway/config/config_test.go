package config

import (
	"testing"
)

func TestCfg_ReadYaml(t *testing.T) {
	c, err := ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	if c.DB.Host != "localhost" {
		t.Errorf("Expected localhost, got %s", c.DB.Host)
	}
	if c.DB.Username != "koosha" {
		t.Errorf("Expected koosha, got %s", c.DB.Username)
	}
	if c.DB.Name != "mock_gt_pq_ns_iop" {
		t.Errorf("Expected mk_pq_ns_iop, got %s", c.DB.Name)
	}
}
