package runner

import (
	"time"

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
			Name:    "default-health-check-path",
			Usage:   "default health check path",
			Value:   "/health",
			EnvVars: []string{"DEFAULT_HEALTH_CHECK_PATH"},
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:    "workers",
			Usage:   "number of workers",
			Value:   1,
			EnvVars: []string{"WORKERS"},
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "check-interval",
			Usage:   "check interval",
			Value:   time.Minute,
			EnvVars: []string{"CHECK_INTERVAL"},
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
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "sources-store-dir",
			Usage:   "sources store directory",
			Value:   "/sources",
			EnvVars: []string{"SOURCES_STORE_DIR"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "sources-store-file",
			Usage:   "sources store file",
			Value:   "sources.yaml",
			EnvVars: []string{"SOURCES_STORE_FILE"},
		}),
		altsrc.NewUint64Flag(&cli.Uint64Flag{
			Name:    "transient-errors-max-retries",
			Usage:   "transient errors max retries",
			Value:   3,
			EnvVars: []string{"TRANSIENT_ERRORS_MAX_RETRIES"},
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "transient-errors-retry-wait",
			Usage:   "transient errors retry wait",
			Value:   time.Second,
			EnvVars: []string{"TRANSIENT_ERRORS_RETRY_WAIT"},
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "http-client-timeout",
			Usage:   "http client timeout",
			Value:   time.Second * 5,
			EnvVars: []string{"HTTP_CLIENT_TIMEOUT"},
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "acceptable-response-time",
			Usage:   "acceptable response time",
			Value:   time.Second * 3,
			EnvVars: []string{"ACCEPTABLE_RESPONSE_TIME"},
		}),
	}
}
