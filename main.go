package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/e10k/dbdl/commands"
	"github.com/e10k/dbdl/settings"
)

func main() {
	command, port := parseInput(os.Args)

	if command == nil {
		fmt.Println("You need to specify a command.")
		os.Exit(1)
	}

	settings := settings.NewSettings()
	if port != nil {
		settings.Port = *port
	}

	switch *command {
	case "install":
		commands.HandleInstallCommand(settings)
	case "serve":
		commands.HandleServeCommand(settings)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		os.Exit(1)
	}
}

func parseInput(args []string) (*string, *int) {
	var command string
	var port int

	args = args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	// the first argument is considered to be the command, as long as it's not actually a flag
	if !strings.HasPrefix(args[0], "-") {
		command = args[0]
	}

	flag.IntVar(&port, "port", 2222, "The port")
	flag.CommandLine.Parse(os.Args[2:])

	return &command, &port
}
