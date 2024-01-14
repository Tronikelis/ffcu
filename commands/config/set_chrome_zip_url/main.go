package set_chrome_zip_url

import (
	"errors"
	"net/url"

	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func action(fConfig *ffcu.Config) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		u := ctx.Args().First()
		if u == "" {
			return errors.New("provide a url")
		}

		parsedUrl, err := url.Parse(u)
		if err != nil {
			return err
		}

		fConfig.ZippedChromeUrl = parsedUrl.String()

		return nil
	}

}

func Command(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name:   "set-chrome-zip",
		Usage:  "Download the zip which has the latest chrome dir from this url",
		Action: action(fConfig),
	}
}
