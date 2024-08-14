package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/e10k/dbdl/commands"
	"github.com/e10k/dbdl/settings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Please specify a command.")
		os.Exit(1)
	}

	settings := settings.NewSettings()

	command := flag.Arg(0)

	switch command {
	case "install":
		commands.HandleInstallCommand(flag.Args()[1:], settings)
	case "serve":
		commands.HandleServeCommand(flag.Args()[1:], settings)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
