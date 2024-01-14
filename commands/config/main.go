package config

import (
	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func GetCommand(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name: "config",
		Subcommands: []*cli.Command{
			setProfileDir(fConfig),
			setUserJsUrl(fConfig),
			setChromeZipUrl(fConfig),
		},
	}
}
