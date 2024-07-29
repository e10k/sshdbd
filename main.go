package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Connection struct {
	Host     string
	Port     int
	Dbname   string
	Username string
	Password string
}

type Config struct {
	Connections map[string]Connection
}

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

	var conf Config
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

	var gzbuf bytes.Buffer
	gz := gzip.NewWriter(&gzbuf)

	dumpSchemaCmd := exec.Command(
		"mysqldump",
		"--single-transaction",
		"--databases",
		"--no-data",
		"--skip-add-drop-table",
		"--verbose",
		"-h",
		conn.Host,
		fmt.Sprintf("-u%s", conn.Username),
		fmt.Sprintf("-p%s", conn.Password),
		conn.Dbname,
	)
	dumpSchemaCmd.Stdout = gz
	dumpSchemaCmd.Stderr = os.Stderr
	dumpSchemaCmd.Run()

	dumpDataCmd := exec.Command(
		"mysqldump",
		"--single-transaction",
		"--tz-utc=false",
		"--no-create-info",
		"--verbose",
		"-h",
		conn.Host,
		fmt.Sprintf("-u%s", conn.Username),
		fmt.Sprintf("-p%s", conn.Password),
		conn.Dbname,
		"--ignore-table=sakila.film",
	)
	dumpDataCmd.Stdout = gz
	dumpDataCmd.Stderr = os.Stderr
	dumpDataCmd.Run()

	gz.Flush()
	gz.Close()
	io.Copy(os.Stdout, &gzbuf)

	// create a config (read from options)
	// dump db structure; put mysqldump stdout output in a buffer; send stderr output to stderr
	// dump db data; put mysqldump stdout output in a buffer; send stderr output to stderr
	// gzip the buffer while it's written
	// stream the gzipped output to stdout
	// show progress bar while the dump progresses
}
