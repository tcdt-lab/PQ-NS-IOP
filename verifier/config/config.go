package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB       DB       `yaml:"DB"`
	Security Security `yaml:"Security"`
	// Ensure this matches the YAML key
}

type DB struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Security struct {
	LoginType   string `yaml:"login_type"`
	MlDSAScheme string `yaml:"mldsa_scheme"`
	MlKEMScheme string `yaml:"mlkem_scheme"`
}

func ReadYaml() (*Config, error) {
	c := &Config{}

	yamlFile, err := os.ReadFile("/home/koosha/Desktop/Thesis/impl/PQ-NS-IOP/verifier/config/config.yaml")
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {

		zap.L().Error("Error unmarshalling config.yaml file", zap.Error(err))
		return nil, err
	}

	return c, nil
}
