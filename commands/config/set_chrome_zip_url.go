package config

import (
	"errors"
	"net/url"

	"github.com/Tronikelis/ffcu/configuration"
	"github.com/urfave/cli/v2"
)

func setChromeZipUrl(config *configuration.Config) *cli.Command {
	action := func(ctx *cli.Context) error {
		u := ctx.Args().First()
		if u == "" {
			return errors.New("provide a url")
		}

		parsedUrl, err := url.Parse(u)
		if err != nil {
			return err
		}

		config.ZippedChromeUrl = parsedUrl.String()

		return nil
	}

	return &cli.Command{
		Name:        "set-chrome-zip",
		Description: "Download the zip which has the latest chrome dir from this url",
		Action:      action,
	}
}
