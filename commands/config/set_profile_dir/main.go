package set_profile_dir

import (
	"errors"
	"os"
	"path"

	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func action(fConfig *ffcu.Config) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		dir := ctx.Args().First()
		if dir == "" {
			return errors.New("please provide a directory")
		}

		if _, err := os.Stat(path.Join(dir, "prefs.js")); os.IsNotExist(err) {
			return errors.New("the directory you provided does not have a prefs.js")
		}

		fConfig.ProfileDir = dir

		return nil
	}
}

func Command(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name:   "set-profile-dir",
		Usage:  "The dir where the user.js and chrome dir is",
		Action: action(fConfig),
	}
}
