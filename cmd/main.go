package main

import (
	"log"
	"os"

	"github.com/adamkisala/health/cmd/runner"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "health-checker",
		Usage: "run health-checker commands",
		Commands: []*cli.Command{
			runner.NewCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
