package print

import (
	"encoding/json"
	"fmt"

	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/urfave/cli/v2"
)

func action(fConfig *ffcu.Config) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		pretty, err := json.MarshalIndent(*fConfig, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(pretty))

		return nil
	}
}

func Command(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name:   "print",
		Usage:  "Prints the current ffcu config",
		Action: action(fConfig),
	}
}
