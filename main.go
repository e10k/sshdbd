package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/BurntSushi/toml"
	"github.com/e10k/dbdl/config"
	"github.com/e10k/dbdl/server"
)

func main() {
	var conf config.Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}
	// fmt.Fprintf(os.Stderr, "conf: %v\n", conf)

	log.Println("starting ssh server on port 2222...")

	err = server.NewServer(conf).ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
