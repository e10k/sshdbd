package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/connections"
)

type Config struct {
	ConfigDir string
	Connections    connections.Connections
	Port      int
}

func NewConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error reading the home dir: %s", err)
	}

	configDir := home + "/.sshdbd"

	var connections connections.Connections

	return &Config{
		ConfigDir: configDir,
		Connections:    connections,
		Port:      2222,
	}, nil
}

func (s *Config) LoadConnections() error {
	file := s.ConfigDir + "/connections.toml"
	_, err := toml.DecodeFile(file, &s.Connections)
	if err != nil {
		return fmt.Errorf("error decoding file %v: %v", file, err)
	}

	return nil
}
