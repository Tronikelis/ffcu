package commands

import (
	"github.com/Tronikelis/ffcu/commands/config"
	"github.com/Tronikelis/ffcu/commands/update"
	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func Commands(fConfig *ffcu.Config) []*cli.Command {
	return []*cli.Command{
		config.Command(fConfig),
		update.Command(fConfig),
	}
}
