package commands

import (
	"github.com/Tronikelis/ffcu/commands/config"
	"github.com/Tronikelis/ffcu/commands/update"
	"github.com/Tronikelis/ffcu/configuration"
	"github.com/urfave/cli/v2"
)

func GetCommands(configuration *configuration.Config) []*cli.Command {
	return []*cli.Command{
		config.GetCommand(configuration),
		update.GetCommand(configuration),
	}
}
