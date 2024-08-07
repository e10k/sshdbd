package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/db"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	sessionHandler := func(s ssh.Session) {
		connId, dbName, skippedTables, err := parseInput(s.User())
		if err != nil {
			panic(err)
		}

		conn, ok := conf.Connections[connId]
		if !ok {
			s.Stderr().Write([]byte(fmt.Sprintf("Invalid connection id '%v'.\n", connId)))
			return
		}

		log.Printf("connId: %v, dbName: %v", connId, dbName)

		databases, err := db.GetDatabases(&conn)
		log.Printf("databases: %#v", databases)

		if !slices.Contains(databases, dbName) {
			s.Stderr().Write([]byte(fmt.Sprintf("Couldn't find a database named '%v'.\n", dbName)))
			return
		}

		err = db.Dump(&conn, dbName, skippedTables, s, s.Stderr())
		if err != nil {
			panic(err)
		}
	}

	authHandler := func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, k := range getKeys() {
			known, comment, _, _, err := gossh.ParseAuthorizedKey([]byte(k))
			if err != nil {
				log.Printf("encountered invalid public key: %v\n", k)
				continue
			}

			if ssh.KeysEqual(key, known) {
				fmt.Printf("found valid key, having comment %v\n", comment)
				return true
			}
		}

		return false
	}

	server := &ssh.Server{
		Addr:             ":2222",
		Handler:          sessionHandler,
		PublicKeyHandler: authHandler,
	}

	log.Println("starting ssh server on port 2222...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func parseInput(s string) (string, string, []string, error) {
	s = strings.Trim(s, " ")
	slice := strings.Split(s, ":")

	// a valid input will have one of these forms:
	//    main:sakila (connection id + database name)
	//    main:sakila:table_1,table_2,table_3 (connection id + database name + comma separated list of table names)
	l := len(slice)
	if l == 3 {
		var tables []string
		for _, t := range strings.Split(slice[2], ",") {
			t = strings.Trim(t, " ")
			if len(t) > 0 {
				tables = append(tables, t)
			}
		}
		return slice[0], slice[1], tables, nil
	} else if l == 2 {
		return slice[0], slice[1], nil, nil
	}

	return "", "", nil, fmt.Errorf("unexpected input data format: %s", s)
}

func getKeys() []string {
	content, err := os.ReadFile("authorized_keys")
	if err != nil {
		return nil
	}

	var keys []string
	for _, k := range strings.Split(string(content), "\n") {
		k = strings.Trim(k, " ")
		if len(k) > 0 {
			keys = append(keys, k)
		}
	}

	return keys
}
