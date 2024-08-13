package main

import (
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

func setup() {
	// ~/.sshdbd
	// - if ~/sshdbd does not exist, create it
	// - if ~/sshdbd/config.toml does not exist, create it
	// - if ~/sshdbd/authorized_keys does not exist, create it
	// - if ~/sshdbd/hostkey.pem does not exist, create it
}
