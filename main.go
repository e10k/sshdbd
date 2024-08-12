package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/server"
)

func main() {
	// hk, err := generateHostKeyBytes()
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

func generateHostKeyBytes() ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	pemBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}
	return pem.EncodeToMemory(&pemBlock), nil
}
