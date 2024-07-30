package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/mysqldump"
)

var connId *string
var dbName *string
var skippedTables *string

func init() {
	connId = flag.String("connId", "", "The connection id as defined in config.toml")
	dbName = flag.String("dbName", "", "The database name to be cloned")
	skippedTables = flag.String("skippedTables", "", "Comma separated list of database tables to skip")
}

func main() {

	flag.Parse()

	if len(*connId) == 0 {
		panic("need to specify the connection id")
	}
	if len(*dbName) == 0 {
		panic("need to specify the database name")
	}

	fmt.Fprintf(os.Stderr, "connId: %#v, dbName: %#v\n", *connId, *dbName)

	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	conn, ok := conf.Connections[*connId]
	if !ok {
		panic("invalid connection id")
	}

	fmt.Fprintln(os.Stderr, conn)

	err = mysqldump.Dump(&conn, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
}
