package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	_ "github.com/go-sql-driver/mysql"

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
			s.Stderr().Write([]byte(fmt.Sprintf("Invalid connection id '%v'.\n", connId)))
			return
		}

		log.Printf("connId: %v, dbName: %v", connId, dbName)

		databases, err := getDatabases(&conn)
		log.Printf("databases: %#v", databases)

		if !slices.Contains(databases, dbName) {
			s.Stderr().Write([]byte(fmt.Sprintf("Couldn't find a database named '%v'.\n", dbName)))
			return
		}

		err = mysqldump.Dump(&conn, dbName, s, s.Stderr())
		if err != nil {
			panic(err)
		}
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return ssh.KeysEqual(key, known)
	})

	log.Println("starting ssh server on port 2222...")
	err = ssh.ListenAndServe(":2222", nil, publicKeyOption)
	if err != nil {
		log.Fatal(err)
	}
}

func parseInput(s string) (string, string, error) {
	s = strings.Trim(s, " ")
	slice := strings.Split(s, ":")
	if len(slice) != 2 {
		return "", "", fmt.Errorf("unexpected input data format: %s", s)
	}

	return slice[0], slice[1], nil
}

func getDatabases(conn *config.Connection) ([]string, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8mb4&parseTime=True&loc=Local", conn.Username, conn.Password, conn.Host, conn.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("could not connect to the database")
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		return nil, errors.New("could not list the databases")
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, errors.New("failed fetching the databases")
		}
		databases = append(databases, dbName)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("failed fetching the databases")
	}

	return databases, nil
}
