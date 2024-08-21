package main

import (
	"fmt"
	"log"
	"os"

	"github.com/e10k/dbdl/commands"
	"github.com/e10k/dbdl/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}

	app := &cli.App{
		Usage: "A SSH server for downloading database dumps",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "create the configuration directory and the required files",
				Action: func(cCtx *cli.Context) error {
					return commands.HandleInstallCommand(config)
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "start the SSH server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Value:       config.Port,
						Usage:       "listen on port `PORT`",
						Aliases:     []string{"p"},
						Destination: &config.Port,
					},
				},
				Action: func(cCtx *cli.Context) error {
					return commands.HandleServeCommand(config)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
