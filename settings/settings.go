package settings

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
)

type Settings struct {
	ConfigDir string
	Config    config.Config
	Port      int
}

func NewSettings() (*Settings, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error reading the home dir: %s", err)
	}

	configDir := home + "/.sshdbd"

	var conf config.Config

	return &Settings{
		ConfigDir: configDir,
		Config:    conf,
		Port:      2222,
	}, nil
}

func (s *Settings) LoadConfig() error {
	file := s.ConfigDir + "/config.toml"
	_, err := toml.DecodeFile(file, &s.Config)
	if err != nil {
		return fmt.Errorf("error decoding file %v: %v", file, err)
	}

	return nil
}
