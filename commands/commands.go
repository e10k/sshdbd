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

func HandleInstallCommand(args []string, settings *settings.Settings) {
	_ = args

	configDir := settings.ConfigDir

	_, err := os.Stat(configDir)
	if err == nil {
		log.Fatalf("config dir %v already exists", configDir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("unexpected error: %v", err)
	}

	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		log.Fatalf("error creating the dir %s", err)
	}

	authorizedKeysFile := configDir + "/authorized_keys"
	f, err := os.Create(authorizedKeysFile)
	if err != nil {
		log.Fatalf("error creating %s: %s", authorizedKeysFile, err)
	}

	err = f.Chmod(0600)
	if err != nil {
		log.Fatalf("error setting permissions for %s: %s", authorizedKeysFile, err)
	}

	hk, err := server.GenerateHostKeyBytes()
	hostKeyFile := configDir + "/hostkey.pem"
	err = os.WriteFile(hostKeyFile, hk, 0600)
	if err != nil {
		log.Fatalf("error creating %s: %s", hostKeyFile, err)
	}

	configFile := configDir + "/config.toml"
	f, err = os.Create(configFile)
	if err != nil {
		log.Fatalf("error creating %s: %s", configFile, err)
	}
	f.WriteString(fmt.Sprintf("[connections.main]\nhost = %q\nport = %d\nusername = %q\npassword = %q\n\n", "localhost", 3306, "usr", "pass"))
}

func HandleServeCommand(args []string, settings *settings.Settings) {
	_ = args

	err := settings.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	log.Println("starting ssh server on port 2222...")

	srv := server.NewServer(settings)

	hostKeyFile := settings.ConfigDir + "/hostkey.pem"
	privateBytes, err := os.ReadFile(hostKeyFile)
	if err != nil {
		log.Fatal("failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("failed to parse private key: ", err)
	}

	srv.AddHostKey(private)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
