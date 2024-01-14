package config

import (
	"errors"
	"os"
	"path"

	"github.com/Tronikelis/ffcu/configuration"
	"github.com/urfave/cli/v2"
)

func setProfileDir(config *configuration.Config) *cli.Command {
	action := func(ctx *cli.Context) error {
		dir := ctx.Args().First()
		if dir == "" {
			return errors.New("please provide a directory")
		}

		if _, err := os.Stat(path.Join(dir, "prefs.js")); os.IsNotExist(err) {
			return errors.New("the directory you provided does not have a prefs.js")
		}

		config.ProfileDir = dir

		return nil
	}

	return &cli.Command{
		Name:        "set-profile-dir",
		Description: "The dir where the user.js and chrome dir is",
		Action:      action,
	}
}
