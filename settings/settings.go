package settings

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
)

type Settings struct {
	ConfigDir string
	Config    config.Config
}

func NewSettings() *Settings {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error obtaining the home dir: %s", err)
	}

	configDir := home + "/.sshdbd"

	var conf config.Config

	return &Settings{
		ConfigDir: configDir,
		Config:    conf,
	}
}

func (s *Settings) LoadConfig() error {
	_, err := toml.DecodeFile(s.ConfigDir+"/config.toml", &s.Config)
	if err != nil {
		return err
	}

	return nil
}
