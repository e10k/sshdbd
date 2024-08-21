package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/e10k/dbdl/server"
	"github.com/e10k/dbdl/settings"
	"golang.org/x/crypto/ssh"
)

func HandleInstallCommand(settings *settings.Settings) error {
	configDir := settings.ConfigDir

	_, err := os.Stat(configDir)
	if err == nil {
		return fmt.Errorf("config dir %v already exists", configDir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("unexpected error: %v", err)
	}

	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return fmt.Errorf("error creating directory %s", err)
	}

	authorizedKeysFile := configDir + "/authorized_keys"
	f, err := os.Create(authorizedKeysFile)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", authorizedKeysFile, err)
	}

	err = f.Chmod(0600)
	if err != nil {
		return fmt.Errorf("error setting permissions for %s: %s", authorizedKeysFile, err)
	}

	hk, err := server.GenerateHostKeyBytes()
	hostKeyFile := configDir + "/hostkey.pem"
	err = os.WriteFile(hostKeyFile, hk, 0600)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", hostKeyFile, err)
	}

	connectionsFile := configDir + "/connections.toml"
	f, err = os.Create(connectionsFile)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", connectionsFile, err)
	}
	f.WriteString(fmt.Sprintf("[main]\nhost = %q\nport = %d\nusername = %q\npassword = %q\n\n", "localhost", 3306, "usr", "pass"))

	return nil
}

func HandleServeCommand(settings *settings.Settings) error {
	err := settings.LoadConnections()
	if err != nil {
		return fmt.Errorf("error loading connections: %v", err)
	}

	log.Printf("Starting SSH server on port %v...\n", settings.Port)

	srv := server.NewServer(settings)

	hostKeyFile := settings.ConfigDir + "/hostkey.pem"
	privateBytes, err := os.ReadFile(hostKeyFile)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	srv.AddHostKey(private)

	err = srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("serve error: %v", err)
	}

	return nil
}
