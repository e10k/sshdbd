package main

import (
	"fmt"
	"log"
	"os"

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

	connId := "main"
	dbName := "sakila"

	var conf config.Config
	_, err = toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	conn, ok := conf.Connections[connId]
	if !ok {
		panic("invalid connection id")
	}

	fmt.Fprintln(os.Stderr, conn)

	conn.Dbname = dbName

	ssh.Handle(func(s ssh.Session) {
		//s.Stderr().Write([]byte(comment))
		log.Println(comment)
		err = mysqldump.Dump(&conn, s, s.Stderr())
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
