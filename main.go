package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/mysqldump"
)

func main() {
	connId := "main"
	dbName := "sakila"

	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
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

	err = mysqldump.Dump(&conn, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
}
