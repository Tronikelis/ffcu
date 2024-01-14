package config

import (
	"github.com/Tronikelis/ffcu/commands/config/print"
	"github.com/Tronikelis/ffcu/commands/config/set_chrome_zip_url"
	"github.com/Tronikelis/ffcu/commands/config/set_profile_dir"
	"github.com/Tronikelis/ffcu/commands/config/set_userjs_url"
	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func Command(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name: "config",
		Subcommands: []*cli.Command{
			set_userjs_url.Command(fConfig),
			set_profile_dir.Command(fConfig),
			set_chrome_zip_url.Command(fConfig),
			print.Command(fConfig),
		},
	}
}
