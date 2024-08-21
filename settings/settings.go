package settings

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/connections"
)

type Settings struct {
	ConfigDir string
	Connections    connections.Connections
	Port      int
}

func NewSettings() (*Settings, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error reading the home dir: %s", err)
	}

	configDir := home + "/.sshdbd"

	var connections connections.Connections

	return &Settings{
		ConfigDir: configDir,
		Connections:    connections,
		Port:      2222,
	}, nil
}

func (s *Settings) LoadConnections() error {
	file := s.ConfigDir + "/connections.toml"
	_, err := toml.DecodeFile(file, &s.Connections)
	if err != nil {
		return fmt.Errorf("error decoding file %v: %v", file, err)
	}

	return nil
}
