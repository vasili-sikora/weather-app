package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	C Config `yaml:"service"`
}

type Provider struct {
	Type string `yaml:"type"`
}

type Location struct {
	Lat  float64 `yaml:"lat"`
	Long float64 `yaml:"long"`
}

type Config struct {
	P Provider `yaml:"provider"`
	L Location `yaml:"location"`
}

func Parse(r io.Reader) (Config, error) {
	var configFile ConfigFile

	if err := yaml.NewDecoder(r).Decode(&configFile); err != nil {
		return Config{}, err
	}

	return configFile.C, nil
}
