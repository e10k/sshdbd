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
		fmt.Println("Please specify a command.")
		os.Exit(1)
	}

	settings, err := settings.NewSettings()
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}
	if port != nil {
		settings.Port = *port
	}

	var e error
	switch *command {
	case "install":
		e = commands.HandleInstallCommand(settings)
	case "serve":
		e = commands.HandleServeCommand(settings)
	default:
		e = fmt.Errorf("Unknown command: %s\n", *command)
	}

	if e != nil {
		fmt.Println(e)
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

	flag.IntVar(&port, "port", 2222, "Port")
	flag.CommandLine.Parse(os.Args[2:])

	return &command, &port
}
