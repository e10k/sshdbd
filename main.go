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
	// hk, err := server.GenerateHostKeyBytes()
	// os.WriteFile("test.pem", hk, 0600)

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

	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	// fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	log.Println("starting ssh server on port 2222...")

	srv := server.NewServer(conf)

	privateBytes, err := os.ReadFile("id_rsa")
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

func handleInstallCommand(args []string) {
	//
}

func handleServeCommand(args []string) {
	//
}
