package server

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "load",
			Usage:   "load configuration from yaml `FILE`",
			EnvVars: []string{"LOAD"},
		},
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "services-hosts",
			Usage:   "services hosts to check",
			EnvVars: []string{"SERVICES_HOSTS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "default-health-path",
			Usage:   "default health path",
			Value:   "/health",
			EnvVars: []string{"DEFAULT_HEALTH_PATH"},
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:    "workers",
			Usage:   "number of workers",
			Value:   1,
			EnvVars: []string{"WORKERS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "log-level",
			Usage:   "log level",
			Value:   "info",
			EnvVars: []string{"LOG_LEVEL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "log-format",
			Usage:   "log format",
			Value:   "json",
			EnvVars: []string{"LOG_FORMAT"},
		}),
	}
}
