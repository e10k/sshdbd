package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/server"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Please specify a command.")
		os.Exit(1)
	}

	command := flag.Arg(0)

	switch command {
	case "install":
		handleInstallCommand(flag.Args()[1:])
	case "serve":
		handleServeCommand(flag.Args()[1:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleInstallCommand(args []string) {
	_ = args

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error obtaining the home dir: %s", err)
	}

	configDir := home + "/.sshdbd"
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

func handleServeCommand(args []string) {
	_ = args

	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	// fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	log.Println("starting ssh server on port 2222...")

	srv := server.NewServer(conf)

	privateBytes, err := os.ReadFile("hostkey.pem")
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
