package config

import (
	"github.com/Tronikelis/ffcu/configuration"
	"github.com/urfave/cli/v2"
)

func GetCommand(config *configuration.Config) *cli.Command {
	return &cli.Command{
		Name: "config",
		Subcommands: []*cli.Command{
			setProfileDir(config),
			setUserJsUrl(config),
			setChromeZipUrl(config),
		},
	}
}
