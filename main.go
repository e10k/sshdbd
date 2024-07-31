package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/mysqldump"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	known, comment, _, _, err := gossh.ParseAuthorizedKey([]byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINMgfBf9NEfplAqXjMdiHCPM0J+f6JVJX4BE2SfEkvPr emi"))
	if err != nil {
		log.Fatal(err)
	}

	var conf config.Config
	_, err = toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	ssh.Handle(func(s ssh.Session) {
		log.Println(comment)
		connId, dbName, err := parseInput(s.User())
		if err != nil {
			panic(err)
		}

		conn, ok := conf.Connections[connId]
		if !ok {
			panic("invalid connection id")
		}

		log.Printf("connId: %v, dbName: %v", connId, dbName)

		err = mysqldump.Dump(&conn, dbName, s, s.Stderr())
		if err != nil {
			panic(err)
		}
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return ssh.KeysEqual(key, known)
		// return true // allow all keys, or use ssh.KeysEqual() to compare against known keys
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, publicKeyOption))
}

func parseInput(s string) (string, string, error) {
	s = strings.Trim(s, " ")
	slice := strings.Split(s, ":")
	if len(slice) != 2 {
		return "", "", fmt.Errorf("unexpected input data format: %s", s)
	}

	return slice[0], slice[1], nil
}
