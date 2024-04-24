package runner

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/adamkisala/weaviate-health/internal/healthcheck"
	"github.com/adamkisala/weaviate-health/internal/sources"
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
			logger := buildLogger(c.String("log-format"))

			sourcesProvider := sources.NewProvider(sources.ProviderParams{
				Loader: sources.NewConfig(sources.ConfigParams{
					SourcesStore:         os.DirFS(c.String("sources-store-dir")),
					SourcesStoreFilePath: c.String("sources-store-file"),
					Logger:               logger,
				}),
			})

			requester := healthcheck.NewHTTPRequester(healthcheck.HTTPRequesterParams{
				Doer: healthcheck.NewRetryableClient(healthcheck.RetryableClientParams{
					MaxRetries: c.Uint64("transient-errors-max-retries"),
					RetryWait:  c.Duration("transient-errors-retry-wait"),
					Client: &http.Client{
						Timeout: c.Duration("http-client-timeout"),
					},
				}),
				DefaultHealthCheckPath: c.String("default-health-check-path"),
			})

			responseProcessors := []healthcheck.ResponseProcessor{
				healthcheck.NewLogProcessor(healthcheck.LogProcessorParams{
					Logger:                 logger,
					AcceptableResponseTime: c.Duration("acceptable-response-time"),
				}),
				healthcheck.NewEmailAlerterProcessor(healthcheck.EmailAlerterProcessorParams{
					Logger: logger,
				}),
			}

			healthCheckRunner := healthcheck.NewRunner(healthcheck.RunnerParams{
				Logger:             logger,
				Requester:          requester,
				ResponseProcessors: responseProcessors,
				SourcesProvider:    sourcesProvider,
				Workers:            c.Int64("workers"),
				CheckInterval:      c.Duration("check-interval"),
			})
			controller := NewController(ControllerParams{
				Runner: healthCheckRunner,
				Logger: logger,
			})

			return controller.Run(c.Context)
		},
	}
}

func buildLogger(logFormat string) *slog.Logger {
	switch logFormat {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stderr, nil))
	default:
		return slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
}
