package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB            DB            `yaml:"DB"`
	Security      Security      `yaml:"Security"`
	Server        Server        `yaml:"Server"`
	User          User          `yaml:"User"`
	BootstrapNode BootstrapNode `yaml:"Bootstrap_Node"`
	Client        Client        `yaml:"Client"`
	ZKP           ZKP           `yaml:"ZKP"`
	// Ensure this matches the YAML key
}

type DB struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Security struct {
	LoginType          string `yaml:"login_type"`
	DSAScheme          string `yaml:"dsa_scheme"`
	KEMScheme          string `yaml:"kem_scheme"`
	CryptographyScheme string `yaml:"cryptography_scheme"`
}

type Server struct {
	Protocol   string `yaml:"protocol"`
	Port       string `yaml:"port"`
	BufferSize int    `yaml:"buffer_size"`
	Ip         string `yaml:"ip"`
}

type Client struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
}
type User struct {
	Password string `yaml:"password"`
}

type BootstrapNode struct {
	Ip        string `yaml:"ip"`
	Port      string `yaml:"port"`
	PubKeySig string `yaml:"pub_key_sig"`
	PubKeyKem string `yaml:"pub_key_kem"`
}

type ZKP struct {
	ProverKeyPath string `yaml:"prover_key_path"`
	WitnessPath   string `yaml:"witness_path"`
}

func ReadYaml() (*Config, error) {
	c := &Config{}

	yamlFile, err := os.ReadFile("/home/koosha/Desktop/PQ-NS-IOP/gateway/config/config.yaml")
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
