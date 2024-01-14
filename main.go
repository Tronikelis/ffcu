package main

import (
	"log"
	"os"
	"path"

	"github.com/Tronikelis/ffcu/commands"
	"github.com/Tronikelis/ffcu/ffcu"

	"github.com/urfave/cli/v2"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.MkdirAll(path.Join(homeDir, ".ffcu"), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	configDir := path.Join(homeDir, ".ffcu/config.json")

	config, err := ffcu.ReadConfig(configDir)
	if err != nil {
		log.Println("can't read config, creating a new one")

		config, err = ffcu.CreateConfig(configDir)
		if err != nil {
			log.Fatalln(err)
		}
	}

	app := &cli.App{
		Name:     "ffcu",
		Usage:    "Helper CLI tool to auto update your firefox config",
		Commands: commands.Commands(&config),
		After: func(ctx *cli.Context) error {
			return config.SaveConfig(configDir)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
