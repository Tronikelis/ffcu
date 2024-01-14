package config

import (
	"errors"
	"net/url"

	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func setUserJsUrl(fConfig *ffcu.Config) *cli.Command {
	action := func(ctx *cli.Context) error {
		u := ctx.Args().First()
		if u == "" {
			return errors.New("provide a url")
		}

		parsedUrl, err := url.Parse(u)
		if err != nil {
			return err
		}

		fConfig.UserJsUrl = parsedUrl.String()

		return nil
	}

	return &cli.Command{
		Name:        "set-userjs-url",
		Description: "Download the user.js from this url",
		Action:      action,
	}
}
