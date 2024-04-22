package runner

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewCommand() *cli.Command {
	flags := Flags()
	return &cli.Command{
		Name:   "runner",
		Usage:  "run health checker",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Action: func(c *cli.Context) error {
			controller := NewController(ControllerParams{})

			return controller.Run()
		},
	}
}
