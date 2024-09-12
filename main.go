package main

import (
	"fmt"
	"log"
	"os"

	"github.com/e10k/sshdbd/commands"
	"github.com/e10k/sshdbd/config"
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
				Usage:   "Creates the configuration directory and the required files",
				Action: func(cCtx *cli.Context) error {
					err := commands.HandleInstallCommand(config)
					if err == nil {
						fmt.Printf("Successfully created the configuration directory and the required files (see %s).\n", config.ConfigDir)
					}

					return err
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Starts the SSH server",
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
